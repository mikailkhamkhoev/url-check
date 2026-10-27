# urlcheck

Small Go CLI that checks a bunch of URLs at the same time and prints the
status code + how long each one took - so this project exists for one
reason: **practicing concurrent HTTP in Go**.

## What it does

- Takes URLs as command-line arguments
- Fetches them concurrently, bounded by `-c`
- Applies an overall timeout and a per-request timeout
- Prints one line per URL: status, duration, URL
- Exits with 1 if any URL failed - useful for scripts and CI

## Setup (once)

```bash
go mod init urlcheck
```

No external dependencies - standard library only.

## Run the tests

```bash
go test ./...
```

## Run it locally

```bash
go run . https://example.com https://golang.org https://httpbin.org/status/503
```

Flags:

- `-c 5` - max concurrent requests (default 5)
- `-timeout 15s` - timeout for the whole run (also Ctrl+C cancels)
- `-req-timeout 5s` - timeout for each individual request

Example output:

```
OK     200  124ms  https://example.com
BAD    503  210ms  https://httpbin.org/status/503
ERROR  https://slow.example  context deadline exceeded

2/3 ok, took 1.4s
```

## How it works

- One goroutine per URL, and a `sync.WaitGroup` so `main` waits for all of them.
- A buffered channel `chan struct{}` of size `-c` is used as a semaphore. Each
  goroutine sends into it before making the request and reads from it after.
  That's what limits how many requests are happening at once.
- The overall run has a `context.WithTimeout`, so if the run takes too long
  everything gets cancelled. Each request also has its own timeout via
  `http.Client{Timeout: ...}`.
- Results are written into a slice that's the same length as the input, using
  the goroutine's index. So we don't need a mutex and we keep the order.

## Files

- `main.go` - flags, context, Ctrl+C handling
- `checker.go` - the goroutines + semaphore
- `report.go` - printing
- `checker_test.go` - a couple of tests
- `go.mod` - module definition (standard library only, no `go.sum`)
- `.gitignore` - keeps binaries out of git
- `LICENSE` - MIT
