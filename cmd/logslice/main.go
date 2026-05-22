// Command logslice extracts time-range segments from structured log files.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/logslice/logslice/internal/filter"
	"github.com/logslice/logslice/internal/output"
	"github.com/logslice/logslice/internal/scanner"
)

func main() {
	var (
		from      = flag.String("from", "", "start time (inclusive)")
		to        = flag.String("to", "", "end time (inclusive)")
		format    = flag.String("format", "", "explicit time format (optional)")
		numbered  = flag.Bool("n", false, "prefix output lines with original line numbers")
		inputFile = flag.String("f", "", "input log file (default: stdin)")
	)
	flag.Parse()

	if *from == "" && *to == "" {
		fmt.Fprintln(os.Stderr, "error: at least one of -from or -to must be specified")
		flag.Usage()
		os.Exit(1)
	}

	f, err := filter.New(*from, *to, *format)
	if err != nil {
		log.Fatalf("filter: %v", err)
	}

	var src *os.File
	if *inputFile != "" {
		src, err = os.Open(*inputFile)
		if err != nil {
			log.Fatalf("open: %v", err)
		}
		defer src.Close()
	} else {
		src = os.Stdin
	}

	fmt := output.FormatRaw
	if *numbered {
		fmt = output.FormatNumbered
	}
	w := output.New(output.Options{Format: fmt})
	defer func() {
		if err := w.Flush(); err != nil {
			log.Printf("flush: %v", err)
		}
		fmt.Fprintf(os.Stderr, "matched %d line(s)\n", w.Count())
	}()

	sc := scanner.New(src)
	for sc.Scan() {
		lineNum, line := sc.LineNumber(), sc.Text()
		if f.Match(line) {
			if err := w.WriteLine(lineNum, line); err != nil {
				log.Fatalf("write: %v", err)
			}
		}
	}
	if err := sc.Err(); err != nil {
		log.Fatalf("scan: %v", err)
	}
}
