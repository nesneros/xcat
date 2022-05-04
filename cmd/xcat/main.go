package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nesneros/xcat/pkg/xcat"
)

var version string = "v0.1.4"

func main() {
	out := bufio.NewWriter(flag.CommandLine.Output())
	defer out.Flush()

	if err := run(os.Args, out, os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(3)
	}
}

func run(args []string, out io.Writer, in io.Reader) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.Parse(os.Args)
	flags.Usage = func() { usage(out, flags) }
	showKind := flags.Bool("kind", false, "Print the detected kind")
	flags.Parse(args[1:])

	bufferedInput := bufio.NewReader(in)
	size := bufferedInput.Size()
	reader := xcat.NewReader(bufferedInput, size)
	kind := reader.Kind()

	if *showKind {
		fmt.Fprintf(out, "%v\n", kind)
		return nil
	}
	_, err := io.Copy(out, reader)
	if err != nil {
		return err
	}
	return nil
}

func usage(out io.Writer, flags *flag.FlagSet) {
	fmt.Fprintf(out, "Version: %v, usage of %s:\n", version, os.Args[0])
	flags.PrintDefaults()
	fmt.Fprintf(out, "\nPossible values for kind: %s\n", strings.Join(core.Kinds[:], ", "))
}
