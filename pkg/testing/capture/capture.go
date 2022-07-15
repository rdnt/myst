package capture

import (
	"io"
	"os"
	"sync"
	"testing"

	"gotest.tools/v3/assert"
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

	assert.Equal(c.t, c.capturing, false)
	c.capturing = true

	c.stdout = os.Stdout
	c.stderr = os.Stderr

	var err error
	c.pr, c.pw, err = os.Pipe()
	assert.NilError(c.t, err)

	os.Stdout = c.pw
	os.Stderr = c.pw
}

func (c *Capture) Stop() string {
	c.mux.Lock()
	defer c.mux.Unlock()

	assert.Equal(c.t, c.capturing, true)

	err := c.pw.Close()
	assert.NilError(c.t, err)

	b, err := io.ReadAll(c.pr)
	assert.NilError(c.t, err)

	err = c.pr.Close()
	assert.NilError(c.t, err)

	os.Stdout = c.stdout
	os.Stderr = c.stderr

	c.capturing = false

	return string(b)
}
