package core

//go:generate stringer -type=kind -trimprefix=kind_ -output=enums_stringer.go

import (
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
)

const defaultBufSize = 1024

type kind int8

// experiments has shown that compressing empty input with gzip then the output is 23
const minSizeForDetection = 22

const (
	kind_error kind = iota
	kind_plain
	kind_gzip
	kind_bzip2
)

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

func (x *XcatReader) init(kind kind) {
	allIn := io.MultiReader(bytes.NewReader(x.buf), x.in)
	switch kind {
	case kind_plain:
		x.output = allIn
	case kind_gzip:
		x.output, x.error = gzip.NewReader(allIn)
	case kind_bzip2:
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
		return kind_plain
	}
	c0, c1 := buf[0], buf[1]
	switch {
	case c0 == 0x1f && c1 == 0x8b:
		err := uncompressGzip(buf)
		if err != nil {
			return kind_plain
		}
		return kind_gzip
	case c0 == 'B' && c1 == 'Z' && bytes.Equal(pi[:], buf[4:10]):
		return kind_bzip2
	default:
		return kind_plain
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
