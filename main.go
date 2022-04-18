package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	core "github.com/nesneros/xcat/xcatcore"
)

//go:embed .version
var Version string

func main() {
	var showKindFlag bool
	flag.BoolVar(&showKindFlag, "kind", false, "Print the detected kind")
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "Version: %v, usage of %s:\n", Version, os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(out, "\nPossible values for kind: %s\n", strings.Join(core.Kinds[:], ", "))
	}
	flag.Parse()
	bufferedInput := bufio.NewReader(os.Stdin)
	size := bufferedInput.Size()
	reader := core.NewReader(bufferedInput, size)
	kind := reader.Kind()
	if showKindFlag {
		fmt.Printf("%v\n", kind)
		return
	}
	writer := bufio.NewWriter(os.Stdout)
	written, err := io.Copy(writer, reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Written: %d bytes, error: %v", written, err)
		os.Exit(1)
	}
	writer.Flush()
}
