package core

import (
	"fmt"
	"io"
)

type ReaderWithBuf struct {
	reader  io.Reader
	readPos int
	buf     []byte
}

func (r *ReaderWithBuf) Rest() ([]byte) {
	return r.buf[r.readPos:]
}

func (r *ReaderWithBuf) Read(p []byte, outOfBuf bool) (n int, err error) {
	if len(r.buf) > r.readPos {
		count := copy(p, r.buf[r.readPos:])
		r.readPos += count
		return count, nil
	} else {
		if outOfBuf {
			return r.reader.Read(p)
		} else {
			return 0, fmt.Errorf("buf is empty")
		}
	}
}

func (r *ReaderWithBuf) Reset() {
	if r.readPos > len(r.buf) {
		panic("Read beyond buffer")
	}
	r.readPos = 0
}

// return err = io.EOF if rd is empty
func (r *ReaderWithBuf) Init(rd io.Reader, bufSize int) error {
	r.reader = rd
	buf := make([]byte, bufSize)
	n, err := io.ReadFull(r.reader, buf)
	r.buf = buf[0:n]
	if err != nil && err != io.ErrUnexpectedEOF {
		return err
	}
	return nil
}
