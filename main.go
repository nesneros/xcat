package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
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

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	if err := run(os.Args, flag.CommandLine.Output(), os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(3)
	}
	select {}
}

func run(args []string, w io.Writer, in io.Reader) error {
	out := bufio.NewWriter(w)
	defer out.Flush()
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.Parse(os.Args)
	flags.Usage = func() { usage(flags) }
	showKind := flags.Bool("kind", false, "Print the detected kind")
	flags.Parse(args[1:])

	bufferedInput := bufio.NewReader(in)
	size := bufferedInput.Size()
	reader, err := xcat.NewReader(bufferedInput, size)
	if err != nil {
		return err
	}
	kind := reader.Kind()

	if *showKind {
		fmt.Fprintf(out, "%v\n", kind)
		return nil
	}
	_, err = io.Copy(out, reader)
	return err
}

func usage(flags *flag.FlagSet) {
	printVersionInfo(flags.Output())
	fmt.Fprintf(flags.Output(), "Usage of %s:\n", os.Args[0])
	flags.PrintDefaults()
	fmt.Fprintf(flags.Output(), "\nPossible values for kind: %s\n", strings.Join(xcat.Kinds[:], ", "))
}

func printVersionInfo(out io.Writer) {
	if version == "" {
		return
	}
	fmt.Fprintf(out, "Version: %s, commit: %s, build timestamp: %s\n", version, commit, buildTimestamp)
}
