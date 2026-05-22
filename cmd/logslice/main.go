package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/logslice/logslice/internal/filter"
	"github.com/logslice/logslice/internal/limiter"
	"github.com/logslice/logslice/internal/output"
	"github.com/logslice/logslice/internal/sampler"
	"github.com/logslice/logslice/internal/stats"
)

func main() {
	var (
		from    = flag.String("from", "", "start timestamp (inclusive)")
		to      = flag.String("to", "", "end timestamp (inclusive)")
		fmt_    = flag.String("format", "", "explicit time format (optional)")
		outFmt  = flag.String("out", "raw", "output format: raw|numbered|count")
		step    = flag.Int("sample", 1, "keep every Nth matched line (>=1)")
		limit   = flag.Int("limit", 0, "max lines to output (0 = unlimited)")
		showStats = flag.Bool("stats", false, "print stats summary to stderr")
	)
	flag.Parse()

	if *from == "" || *to == "" {
		fmt.Fprintln(os.Stderr, "error: --from and --to are required")
		os.Exit(1)
	}

	f, err := filter.New(*from, *to, *fmt_)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	smp, err := sampler.New(*step)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	lim := limiter.New(*limit)

	w, err := output.New(os.Stdout, *outFmt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	st := stats.New()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		st.RecordScanned()
		if !f.Match(line) {
			continue
		}
		if !smp.Keep() {
			continue
		}
		if !lim.Keep() {
			break
		}
		st.RecordMatched()
		if err := w.Write(line); err != nil {
			fmt.Fprintf(os.Stderr, "write error: %v\n", err)
			os.Exit(1)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scan error: %v\n", err)
		os.Exit(1)
	}

	w.Flush()
	st.Finish()

	if *showStats {
		st.Print(os.Stderr)
	}
}
