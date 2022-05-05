package xcat

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

type Reader struct {
	buf      []byte
	in       io.Reader
	mappedRd io.Reader
	kind     kind
}

// Create a new xcat.Reader. If bufSize <= 0 the default buffer size is used
func NewReader(in io.Reader, bufSize int) (*Reader, error) {
	if bufSize <= 0 {
		bufSize = defaultBufSize
	}
	if bufSize < minSizeForDetection {
		bufSize = minSizeForDetection
	}
	buf := make([]byte, bufSize)
	n, err := io.ReadFull(in, buf)
	if err == io.ErrUnexpectedEOF {
		err = nil
		buf = buf[:n]
	}
	if err != nil {
		return nil, err
	}
	kind := detectKind(buf)
	allIn := io.MultiReader(bytes.NewReader(buf), in)
	var output io.Reader
	switch kind {
	case kindPlain:
		output = allIn
	case kindGzip:
		output, err = gzip.NewReader(allIn)
	case kindBzip2:
		output = bzip2.NewReader(allIn)
	default:
		panic(fmt.Sprintf("INTERNAL ERROR: Invalid kind: %v", kind))
	}
	if err != nil {
		return nil, err
	}
	return &Reader{in: in, buf: buf, mappedRd: output, kind: kind}, nil
}

func (x *Reader) Kind() kind {
	return x.kind
}

func (x *Reader) Read(p []byte) (n int, err error) {
	return x.mappedRd.Read(p)
}

func detectKind(buf []byte) kind {
	size := len(buf)
	if size < minSizeForDetection {
		return kindPlain
	}
	c0, c1 := buf[0], buf[1]
	switch {
	case c0 == 0x1f && c1 == 0x8b:
		// magic values are only two bytes. Decompressing a bit to make extra sure it is gzip
		err := decompressGzip(buf)
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

func decompressGzip(in []byte) error {
	inRd := bytes.NewReader(in)
	rd, err := gzip.NewReader(inRd)
	if err != nil {
		return err
	}
	_, err = io.Copy(io.Discard, rd)
	switch err {
	case io.ErrUnexpectedEOF, nil:
		return nil
	default:
		return err
	}
}
