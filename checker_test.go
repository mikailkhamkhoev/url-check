package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckURLs_AllOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer srv.Close()

	client := &http.Client{Timeout: time.Second}
	results := CheckURLs(context.Background(), client, []string{srv.URL, srv.URL, srv.URL}, 3)

	for _, r := range results {
		if r.Err != nil {
			t.Fatalf("unexpected error: %v", r.Err)
		}
		if r.Status != 200 {
			t.Fatalf("expected 200, got %d", r.Status)
		}
	}
}

func TestCheckURLs_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	client := &http.Client{Timeout: time.Second}
	results := CheckURLs(ctx, client, []string{srv.URL}, 2)

	if results[0].Err == nil {
		t.Fatalf("expected a timeout error, got %+v", results[0])
	}
}
