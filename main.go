package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"os"
	"strings"

	flag "github.com/spf13/pflag"

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
	if err := run(os.Args, os.Stdout, os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(3)
	}
}

func run(args []string, w io.Writer, in io.Reader) error {
	out := bufio.NewWriter(w)
	defer out.Flush()
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.Usage = func() { usage(flags, os.Stderr) }
	showKind := flags.BoolP("kind", "k", false, "Print the detected kind")
	showLicense := flags.BoolP("license", "L", false, "Show license and exit without reading stdin")
	showHelp := flags.BoolP("help", "h", false, "Show help and exit without reading stdin")
	flags.Parse(args[1:])

	switch {
	case *showLicense:
		fmt.Fprintf(out, "%s\n", license)
		return nil
	case *showHelp:
		usage(flags, os.Stdout)
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

func usage(flags *flag.FlagSet, w io.Writer) {
	fmt.Fprintf(w, "usage: %s [options]\n\n", os.Args[0])
	flags.PrintDefaults()
	fmt.Fprintf(w, "\nPossible values for kind: %s\n", strings.Join(xcat.Kinds[:], ", "))
	printVersionInfo(w)
}

func printVersionInfo(w io.Writer) {
	if version == "" {
		return
	}
	fmt.Fprintf(w, "\nVersion: %s, commit: %s, build timestamp: %s\n", version, commit, buildTimestamp)
}
