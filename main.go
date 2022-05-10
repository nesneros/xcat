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
	if err := bufferAndRun(os.Args, os.Stdout, os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func bufferAndRun(args []string, w io.Writer, r io.Reader) error {
	writer := bufio.NewWriter(w)
	defer writer.Flush()
	return run(args, writer, bufio.NewReader(r))
}

func run(args []string, w io.Writer, in io.Reader) error {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	var printKind, showLicense, showHelp bool
	flags.BoolVarP(&printKind, "kind", "k", false, "Print the detected kind")
	flags.BoolVarP(&showLicense, "license", "L", false, "Show license and exit without reading stdin")
	flags.BoolVarP(&showHelp, "help", "h", false, "Show help and exit without reading stdin")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	restArgs := flags.Args()
	if len(restArgs) > 0 {
		return fmt.Errorf("Too many args: %v", restArgs)
	}
	if showLicense || showHelp {
		switch {
		case showHelp:
			printHelp(w, flags)
		case showLicense:
			fmt.Fprint(w, license)
		}
		return nil
	}
	reader, err := xcat.NewReader(in, 2048)
	if err != nil {
		return err
	}
	kind := reader.Kind()

	if printKind {
		fmt.Fprintf(w, "%s\n", kind)
		return nil
	}
	_, err = io.Copy(w, reader)
	return err
}

func printHelp(w io.Writer, flags *flag.FlagSet) {
	fmt.Fprintf(w, "usage: %s [options]\n\n%s", os.Args[0], flags.FlagUsages())
	fmt.Fprintf(w, "\nPossible values for kind: %s\n", strings.Join(xcat.Kinds[:], ", "))
	printVersionInfo(w)
}

func printVersionInfo(w io.Writer) {
	if version == "" {
		return
	}
	fmt.Fprintf(w, "\nVersion: %s, commit: %s, build timestamp: %s\n", version, commit, buildTimestamp)
}
