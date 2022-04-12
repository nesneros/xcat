package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/nesneros/xcat/core"
)

func main() {
	bufferedInput := bufio.NewReader(os.Stdin)
	size := bufferedInput.Size()
	reader := core.NewReader(bufferedInput, size)
	writer := bufio.NewWriter(os.Stdout)
	written, err := io.Copy(writer, reader)
	if err == nil {
		writer.Flush()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Written: %d, err: %v", written, err)
	}
}
