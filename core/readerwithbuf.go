package core

import (
	"fmt"
	"io"
)

type readerWithBuf struct {
	reader  io.Reader
	readPos int
	buf     []byte
}

func (r *readerWithBuf) rest() []byte {
	return r.buf[r.readPos:]
}

func (r *readerWithBuf) Read(p []byte) (n int, err error) {
	return r.read(p, true)
}

func (r *readerWithBuf) read(p []byte, outOfBuf bool) (n int, err error) {
	if len(r.buf) > r.readPos {
		count := copy(p, r.buf[r.readPos:])
		r.readPos += count
		return count, nil
	}
	if outOfBuf {
		return r.reader.Read(p)
	}
	return 0, fmt.Errorf("buf is empty")
}

func (r *readerWithBuf) reset() {
	if r.readPos > len(r.buf) {
		panic("Read beyond buffer")
	}
	r.readPos = 0
}

// return err = io.EOF if rd is empty
func (r *readerWithBuf) init(rd io.Reader, bufSize int) error {
	r.reader = rd
	buf := make([]byte, bufSize)
	n, err := io.ReadFull(r.reader, buf)
	r.buf = buf[0:n]
	if err != nil && err != io.ErrUnexpectedEOF {
		return err
	}
	return nil
}
