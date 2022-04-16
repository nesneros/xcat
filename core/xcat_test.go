package core

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func plainReader(s string) io.Reader {
	return strings.NewReader(s)
}

func gzipToBytes(s string) []byte {
	var b bytes.Buffer
	w1 := bufio.NewWriter(&b)
	w2 := gzip.NewWriter(w1)
	w2.Write([]byte(s))
	w2.Close()
	w1.Flush()
	return b.Bytes()
}

func gzipReader(s string) io.Reader {
	return bytes.NewReader(gzipToBytes(s))
}

func TestDetectGzip(t *testing.T) {
	assert := assert.New(t)
	bb := gzipToBytes("Hello world")
	// cut one byte away
	kind, uncompressed := detectKind(bb[:len(bb)-1])
	assert.Equal(kind_gzip, kind)
	assert.Nil(uncompressed)
	// change a byte (i.e. no valid gzip header)
	bb[len(bb)-1]++
	kind, uncompressed = detectKind(bb)
	assert.Equal(kind_plain, kind)
	assert.Nil(uncompressed)
	// change a byte (i.e. no valid gzip header)
	bb[len(bb)-2]--
	kind, uncompressed = detectKind(bb)
	assert.Equal(kind_plain, kind)
	assert.Nil(uncompressed)
}

func TestPlain(t *testing.T) {
	assert := assert.New(t)
	rd := plainReader("abc")
	xcatRd := NewReader(rd, 100)
	out, e := io.ReadAll(xcatRd)
	assert.NoError(e)
	assert.Equal("abc", string(out))
}

func TestGzip(t *testing.T) {
	assert := assert.New(t)
	rd := gzipReader("abc")
	xcatRd := NewReader(rd, 100)
	out, e := io.ReadAll(xcatRd)
	assert.NoError(e)
	assert.Equal("abc", string(out))
}
