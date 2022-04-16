package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/nesneros/xcat/core"
)

func main() {
	var kindFlag bool
	flag.BoolVar(&kindFlag, "kind", false, "Write the detected kind to stderr")
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(out, "\nPossible value for kind: ")
		for i, e := range core.Kinds {
			if i > 0 {
				fmt.Fprint(out, ", ")
			}
			fmt.Fprint(out, e)
		}
		fmt.Fprintln(out)
	}
	flag.Parse()
	bufferedInput := bufio.NewReader(os.Stdin)
	size := bufferedInput.Size()
	reader := core.NewReader(bufferedInput, size)
	kind := reader.Kind()
	if kindFlag {
		fmt.Fprintf(os.Stderr, "%v\n", kind)
	}

	writer := bufio.NewWriter(os.Stdout)
	written, err := io.Copy(writer, reader)
	if err == nil {
		writer.Flush()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Written: %d, err: %v", written, err)
	}
}
