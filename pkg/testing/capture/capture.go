package capture

import (
	"io"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type Capture struct {
	mux       sync.Mutex
	t         *testing.T
	capturing bool
	stdout    *os.File
	stderr    *os.File
	pr        *os.File
	pw        *os.File
}

func New(t *testing.T) *Capture {
	return &Capture{t: t}
}

func (c *Capture) Start() {
	c.mux.Lock()
	defer c.mux.Unlock()

	require.Equal(c.t, c.capturing, false)
	c.capturing = true

	c.stdout = os.Stdout
	c.stderr = os.Stderr

	var err error
	c.pr, c.pw, err = os.Pipe()
	require.NoError(c.t, err)

	os.Stdout = c.pw
	os.Stderr = c.pw
}

func (c *Capture) Stop() string {
	c.mux.Lock()
	defer c.mux.Unlock()

	require.Equal(c.t, c.capturing, true)

	err := c.pw.Close()
	require.NoError(c.t, err)

	b, err := io.ReadAll(c.pr)
	require.NoError(c.t, err)

	err = c.pr.Close()
	require.NoError(c.t, err)

	os.Stdout = c.stdout
	os.Stderr = c.stderr

	c.capturing = false

	return string(b)
}
