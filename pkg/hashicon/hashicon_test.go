package hashicon

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestNew(t *testing.T) {
	b := make([]byte, 16)
	_, err := New(b)
	assert.Assert(t, err != nil)

	b = make([]byte, 48)
	_, err = New(b)
	assert.Assert(t, err != nil)

	b = make([]byte, 16384)
	_, err = New(b)
	assert.Assert(t, err != nil)

	b = make([]byte, 8192)
	_, err = New(b)
	assert.Equal(t, err, nil)

	b = make([]byte, 32)
	h, err := New(b)
	assert.Equal(t, err, nil)
	assert.Equal(t, h.Stride, 8)
	assert.Equal(t, len(h.Pix), 64)
}
