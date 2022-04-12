package core

import (
	"io"
)

const defaultBufSize = 8 * 1024

type state int8

const (
	initial = state(iota)
	readBuffer
	readStream
)

type kind int8

const (
	kind_error = kind(iota)
	kind_plain
	kind_gzip
)

type XcatReader struct {
	input ReaderWithBuf
	state state
	kind  kind
	error error
}

func NewReader(in io.Reader, bufSize int) *XcatReader {
	if bufSize <= 0 {
		bufSize = defaultBufSize
	}
	buf := make([]byte, 0, bufSize)
	_, err := io.ReadFull(in, buf)
	var k kind
	if err == nil {
		//kind = detectKind(buf)
	}
	result := XcatReader{state: readBuffer, kind: k, error: err}
	result.input.Init(in, bufSize)
	return &result
}

func (x *XcatReader) Read(p []byte) (n int, err error) {
	return x.input.Read(p, true)
}

func detectKind(buf []byte) (kind kind, uncompressed []byte) {
	size := len(buf)
	if size < 64 {
		return kind_plain, nil
	}
	c0, c1 := buf[0], buf[1]
	switch {
	case c0 == 0x1f && c1 == 0x8b:
		return kind_gzip, nil
	default:
		return kind_plain, nil
	}
}
