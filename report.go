package main

import (
	"fmt"
	"io"
	"time"
)

type Summary struct {
	Total int
	OK    int
}

func WriteReport(w io.Writer, results []Result, elapsed time.Duration) Summary {
	summary := Summary{Total: len(results)}

	for _, r := range results {
		if r.Err != nil {
			fmt.Fprintf(w, "%-5s  %s  %v\n", "ERROR", r.URL, r.Err)
			continue
		}

		status := "OK"
		if r.Status >= 200 && r.Status < 400 {
			summary.OK++
		} else {
			status = "BAD"
		}

		fmt.Fprintf(w, "%-5s  %d  %v  %s\n",
			status, r.Status, r.Duration.Round(time.Millisecond), r.URL)
	}

	fmt.Fprintf(w, "\n%d/%d ok, took %v\n",
		summary.OK, summary.Total, elapsed.Round(time.Millisecond))
	return summary
}
