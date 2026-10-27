package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	maxConc := flag.Int("c", 5, "max concurrent requests")
	totalTimeout := flag.Duration("timeout", 15*time.Second, "timeout for the whole run")
	reqTimeout := flag.Duration("req-timeout", 5*time.Second, "timeout for each request")
	flag.Parse()

	urls := flag.Args()
	if len(urls) == 0 {
		fmt.Fprintln(os.Stderr, "usage: urlcheck [-c 5] [-timeout 15s] [-req-timeout 5s] url1 url2 ...")
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *totalTimeout)
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		fmt.Println("\ninterrupted, stopping...")
		cancel()
	}()

	client := &http.Client{Timeout: *reqTimeout}

	start := time.Now()
	results := CheckURLs(ctx, client, urls, *maxConc)
	elapsed := time.Since(start)

	summary := WriteReport(os.Stdout, results, elapsed)

	if summary.OK < summary.Total {
		os.Exit(1)
	}
}
