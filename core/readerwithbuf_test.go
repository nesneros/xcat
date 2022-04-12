package core_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"

	xcat "github.com/nesneros/xcat/core"
)

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func FuzzFoo(f *testing.F) {
	f.Add("abc")
	f.Fuzz(func(t *testing.T, str string) {
		test(t, str)
	})
}

func TestX(t *testing.T) {
	test(t, "abc")
}

func test(t *testing.T, str string) {
	if len(str) == 0 {
		return
	}
	var lenstr int = len(str)
	bufLength := min(lenstr, 10)
	assert := assert.New(t)
	fmt.Fprintf(os.Stderr, "Input: %s\n", str)

	r3 := strings.NewReader(str)
	r := xcat.ReaderWithBuf{}
	r.Init(r3, bufLength)

	assert.Equal(str, string(r.Rest()))
	b20 := make([]byte, 20)
	b2 := make([]byte, 2)
	n, err := r.Read(b20, false)

	assert.NoError(err)
	assert.Equal(len(str), n)
	assert.Equal(string(r.Rest()), "")

	if len(str) <= 2 {
		return
	}
	r.Reset()
	n, err = r.Read(b2, false)
	assert.NoError(err)
	assert.Equal(2, n)
	fmt.Fprintf(os.Stderr, "Input2: %s\n", str)
	tmp := str[2:]
	assert.Equal(string(r.Rest()), tmp)
}
