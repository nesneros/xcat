package xcatcore

//go:generate stringer -type=kind -trimprefix=kind -output=enums_stringer.go

import (
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
)

const (
	defaultBufSize      = 1024
	minSizeForDetection = 18
)

type kind int8

const (
	kindPlain kind = iota
	kindGzip
	kindBzip2
)

// All possible values for kind
var Kinds [kindBzip2 + 1]string

func init() {
	for i := range Kinds {
		// this works regardless if stringer has been generated or not
		Kinds[i] = fmt.Sprintf("%v", kind(i))
	}
}

type XcatReader struct {
	buf    []byte
	in     io.Reader
	output io.Reader
	kind   kind
	error  error
}

func NewReader(in io.Reader, bufSize int) *XcatReader {
	if bufSize <= 0 {
		bufSize = defaultBufSize
	}
	if bufSize < minSizeForDetection {
		bufSize = minSizeForDetection
	}
	buf := make([]byte, bufSize)
	n, err := io.ReadFull(in, buf)
	var k kind
	if err == io.ErrUnexpectedEOF {
		err = nil
		buf = buf[:n]
	}
	result := XcatReader{in: in, buf: buf, error: err}
	if err == nil {
		k = detectKind(buf)
		result.init(k)
	}
	return &result
}

func (x *XcatReader) Kind() kind {
	return x.kind
}

func (x *XcatReader) init(kind kind) {
	x.kind = kind
	allIn := io.MultiReader(bytes.NewReader(x.buf), x.in)
	switch kind {
	case kindPlain:
		x.output = allIn
	case kindGzip:
		x.output, x.error = gzip.NewReader(allIn)
	case kindBzip2:
		x.output = bzip2.NewReader(allIn)
	default:
		panic(fmt.Sprintf("Invalid kind: %v", kind))
	}
}

func (x *XcatReader) Read(p []byte) (n int, err error) {
	if x.error != nil {
		return 0, x.error
	}
	return x.output.Read(p)
}

func detectKind(buf []byte) kind {
	size := len(buf)
	if size < minSizeForDetection {
		return kindPlain
	}
	c0, c1 := buf[0], buf[1]
	switch {
	case c0 == 0x1f && c1 == 0x8b:
		err := uncompressGzip(buf)
		if err != nil {
			return kindPlain
		}
		return kindGzip
	case c0 == 'B' && c1 == 'Z' && bytes.Equal(pi[:], buf[4:10]):
		return kindBzip2
	default:
		return kindPlain
	}
}

var pi = [...]byte{0x31, 0x41, 0x59, 0x26, 0x53, 0x59}

func uncompressGzip(in []byte) error {
	inRd := bytes.NewReader(in)
	var tmp any = inRd
	_, ok := tmp.(io.ByteReader)
	if !ok {
		panic("internal error")
	}
	rd, err1 := gzip.NewReader(inRd)
	if err1 != nil {
		return err1
	}
	_, err2 := io.Copy(io.Discard, rd)
	switch err2 {
	case io.ErrUnexpectedEOF:
		return nil
	case nil:
		return nil
	default:
		return err2
	}
}
