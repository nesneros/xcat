package core

// disabled go:generate stringer -type=kind -trimprefix=kind_ -output=enums_stringer.go

import (
	"bytes"
	"compress/gzip"
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
)

type XcatReader struct {
	input  readerWithBuf
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
	result := XcatReader{error: err}
	if err == nil {
		k, _ = detectKind(buf)
		result.init(in, buf, k)
	}
	return &result
}

func (x *XcatReader) init(in io.Reader, buf []byte, kind kind) {
	x.input.reader = in
	x.input.buf = buf
	switch kind {
	case kind_plain:
		x.output = &x.input
	case kind_gzip:
		x.output, x.error = gzip.NewReader(&x.input)
	}
}

func (x *XcatReader) Read(p []byte) (n int, err error) {
	if x.error != nil {
		return 0, x.error
	}
	return x.output.Read(p)
}

func detectKind(buf []byte) (kind kind, uncompressed []byte) {
	size := len(buf)
	if size < minSizeForDetection {
		return kind_plain, nil
	}
	c0, c1 := buf[0], buf[1]
	switch {
	case c0 == 0x1f && c1 == 0x8b:
		uncompressed, err := uncompressGzip(buf)
		if err != nil {
			return kind_plain, nil
		}
		return kind_gzip, uncompressed
	default:
		return kind_plain, nil
	}
}

func uncompressGzip(in []byte) (out []byte, err1 error) {
	inRd := bytes.NewReader(in)
	var tmp any = inRd
	_, ok := tmp.(io.ByteReader)
	if !ok {
		panic("internal error")
	}
	rd, err1 := gzip.NewReader(inRd)
	if err1 != nil {
		return nil, err1
	}
	uncompressed, err2 := io.ReadAll(rd)
	switch err2 {
	case io.ErrUnexpectedEOF:
		return nil, nil
	case nil:
		return uncompressed, nil
	default:
		return nil, err2
	}
}
