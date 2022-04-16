package core

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
)

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func TestReadOnlyBuffer(t *testing.T) {
	assert := assert.New(t)

	r3 := strings.NewReader("abc")
	r := readerWithBuf{}
	r.init(r3, 10)

	assert.Equal("abc", string(r.rest()))
	b10 := make([]byte, 10)
	b2 := make([]byte, 2)
	n, err := r.read(b10, false)

	assert.NoError(err)
	assert.Equal(3, n)
	assert.Equal(len(r.rest()), 0)

	r.reset()
	n, err = r.read(b2, false)
	assert.NoError(err)
	assert.Equal(2, n)
	assert.Equal(string(r.rest()), "c")
	n, err = r.read(b2, true)
	assert.NoError(err)
	assert.Equal(1, n)
	assert.Equal("", string(r.rest()))
	assert.Equal("c", string(b2[:n]))

	_, err = r.read(b2, true)
	assert.Equal(io.EOF, err)
}

func TestReadFromStream(t *testing.T) {
	assert := assert.New(t)

	r3 := strings.NewReader("0123456789HELLO")
	r := readerWithBuf{}
	r.init(r3, 10)

	assert.Equal("0123456789", string(r.rest()))
	b10 := make([]byte, 10)
	n, err := r.read(b10, false)

	assert.NoError(err)
	assert.Equal(10, n)
	assert.Equal(len(r.rest()), 0)

	n, err = r.read(b10, true)
	assert.Equal(5, n)
	assert.NoError(err)
}
