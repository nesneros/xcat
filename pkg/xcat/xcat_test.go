package xcat

import (
	"bufio"
	"bytes"
	"compress/gzip"
	_ "embed"
	"io"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func gzipToBytes(s string) []byte {
	var b bytes.Buffer
	w1 := bufio.NewWriter(&b)
	w2 := gzip.NewWriter(w1)
	w2.Write([]byte(s))
	w2.Close()
	w1.Flush()
	return b.Bytes()
}

func writeRandom(w io.Writer, n int) (checksum int) {
	var next byte
	for i := 0; i < n; i++ {
		if rand.Intn(1000) == 0 {
			next = byte(rand.Int())
		}
		checksum += int(next)
		b := []byte{next}
		w.Write(b)
	}
	return
}
func TestRandom(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	w1 := bufio.NewWriter(&b)
	w2 := gzip.NewWriter(w1)
	wChecksum := writeRandom(w2, 100000)
	w2.Close()
	w1.Flush()
	cb := b.Bytes()
	xcat, _ := NewReader(bytes.NewReader(cb), 200)
	buf, err := io.ReadAll(xcat)
	assert.NoError(err)
	rChecksum := 0
	for _, e := range buf {
		rChecksum += int(e)
	}
	assert.Equal(wChecksum, rChecksum)
}

func TestDetectGzip(t *testing.T) {
	assert := assert.New(t)
	bb := gzipToBytes("Hello world")
	// cut one byte away
	kind := detectKind(bb[:len(bb)-1])
	assert.Equal(kindGzip, kind)
	// change a byte (i.e. no valid gzip header)
	bb[len(bb)-1]++
	kind = detectKind(bb)
	assert.Equal(kindPlain, kind)
	// change a byte (i.e. no valid gzip header)
	bb[len(bb)-2]--
	kind = detectKind(bb)
	assert.Equal(kindPlain, kind)
}

//go:embed helloworld.bz
var bzHelloWorld string

func TestDetectBzip2(t *testing.T) {
	assert := assert.New(t)
	b := []byte(bzHelloWorld)
	kind := detectKind(b)
	assert.Equal(kindBzip2, kind)
}

func TestPlain(t *testing.T) {
	assert := assert.New(t)
	rd := strings.NewReader("abc")
	xcatRd, err := NewReader(rd, 100)
	assert.NoError(err)
	out, err := io.ReadAll(xcatRd)
	assert.NoError(err)
	assert.Equal("abc", string(out))
}

func TestEmpty(t *testing.T) {
	assert := assert.New(t)
	rd := strings.NewReader("")
	xcatRd, err := NewReader(rd, 100)
	assert.NoError(err)
	out, err := io.ReadAll(xcatRd)
	assert.NoError(err)
	assert.Equal("", string(out))
}
func TestGzip(t *testing.T) {
	testGzip(t, "abc")
	testGzip(t, "\n")
	testGzip(t, "")
}

func testGzip(t *testing.T, s string) {
	assert := assert.New(t)
	rd := bytes.NewReader(gzipToBytes(s))
	xcatRd, err := NewReader(rd, 100)
	assert.NoError(err)
	out, err := io.ReadAll(xcatRd)
	assert.NoError(err)
	assert.Equal(s, string(out))
}

func TestBzip2(t *testing.T) {
	assert := assert.New(t)
	rd := strings.NewReader(bzHelloWorld)
	xcatRd, err := NewReader(rd, 100)
	assert.NoError(err)
	out, err := io.ReadAll(xcatRd)
	assert.NoError(err)
	assert.Equal("Hello World\n", string(out))
}
