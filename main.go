package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nesneros/xcat/pkg/xcat"
)

// nolint: gochecknoglobals
var (
	version        string
	commit         string
	buildTimestamp string
)

//go:embed LICENSE
var license string

func main() {
	if err := run(os.Args, flag.CommandLine.Output(), os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(3)
	}
}

func run(args []string, w io.Writer, in io.Reader) error {
	out := bufio.NewWriter(w)
	defer out.Flush()
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.Usage = func() { usage(flags) }
	showKind := flags.Bool("kind", false, "Print the detected kind")
	showLicense := flags.Bool("license", false, "Show license")
	flags.Parse(args[1:])

	if *showLicense {
		fmt.Fprintf(out, "%s\n", license)
		return nil
	}

	bufferedInput := bufio.NewReader(in)
	size := bufferedInput.Size()
	reader, err := xcat.NewReader(bufferedInput, size)
	if err != nil {
		return err
	}
	kind := reader.Kind()

	if *showKind {
		fmt.Fprintf(out, "%s\n", kind)
		return nil
	}
	_, err = io.Copy(out, reader)
	return err
}

func usage(flags *flag.FlagSet) {
	fmt.Fprintf(flags.Output(), "Usage of %s:\n", os.Args[0])
	flags.PrintDefaults()
	fmt.Fprintf(flags.Output(), "\nPossible values for kind: %s\n", strings.Join(xcat.Kinds[:], ", "))
	printVersionInfo(flags.Output())
}

func printVersionInfo(out io.Writer) {
	if version == "" {
		return
	}
	fmt.Fprintf(out, "\nVersion: %s, commit: %s, build timestamp: %s\n", version, commit, buildTimestamp)
}
