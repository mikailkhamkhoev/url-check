package main

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	URL      string
	Status   int
	Duration time.Duration
	Err      error
}

func CheckURLs(ctx context.Context, client *http.Client, urls []string, maxConc int) []Result {
	results := make([]Result, len(urls))

	sem := make(chan struct{}, maxConc)
	var wg sync.WaitGroup

	for i, url := range urls {
		wg.Add(1)
		go func(i int, url string) {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				results[i] = Result{URL: url, Err: ctx.Err()}
				return
			}
			defer func() { <-sem }()

			results[i] = fetchOne(ctx, client, url)
		}(i, url)
	}

	wg.Wait()
	return results
}

func fetchOne(ctx context.Context, client *http.Client, url string) Result {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return Result{URL: url, Duration: time.Since(start), Err: err}
	}

	resp, err := client.Do(req)
	if err != nil {
		return Result{URL: url, Duration: time.Since(start), Err: err}
	}
	defer resp.Body.Close()

	return Result{
		URL:      url,
		Status:   resp.StatusCode,
		Duration: time.Since(start),
	}
}
