package downloader

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var opts_default Options = Options{DeleteOnError: true}
var opts_keep_on_err Options = Options{DeleteOnError: false}

func TestDownload_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}))
	defer server.Close()

	dir := t.TempDir()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	ctx := context.Background()

	err := Download(opts_default, ctx, client, server.URL+"/file.txt", dir, "file.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "/file.txt"))
	if err != nil {
		t.Fatalf("unable to read file: %v", err)
	}

	if string(data) != "hello" {
		t.Fatalf("unexpected content: %s", string(data))
	}
}

func TestDownload_InvalidURL(t *testing.T) {
	dir := t.TempDir()
	client := &http.Client{}
	ctx := context.Background()

	err := Download(opts_default, ctx, client, "://bad-url", dir, "fallback.txt")
	if err == nil {
		t.Fatalf("expected error for invalid URL")
	}
}

func TestDownload_ContextCancelledImmediately(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.Write([]byte("hello"))
	}))
	defer server.Close()

	dir := t.TempDir()
	client := &http.Client{}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := Download(opts_default, ctx, client, server.URL+"/file.txt", dir, "file.txt")
	if err == nil {
		t.Fatalf("expected context cancellation error")
	}
}

func TestDownload_ContextCancelledDuringTransfer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < 10; i++ {
			w.Write([]byte("chunk"))
			w.(http.Flusher).Flush()
			time.Sleep(500 * time.Millisecond)
		}
	}))
	defer server.Close()

	dir := t.TempDir()
	client := &http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := Download(opts_default, ctx, client, server.URL+"/file.txt", dir, "file.txt")
	if err == nil {
		t.Fatalf("expected cancellation error")
	}
}

func TestDownload_FileStatusOnError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < 10; i++ {
			w.Write([]byte("chunk"))
			w.(http.Flusher).Flush()
			time.Sleep(200 * time.Millisecond)
		}
	}))
	defer server.Close()

	tests := []struct {
		name        string
		opt         Options
		shouldExist bool
	}{
		{
			name:        "delete_on_error",
			opt:         opts_default,
			shouldExist: false,
		},
		{
			name:        "keep_on_error",
			opt:         opts_keep_on_err,
			shouldExist: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			client := &http.Client{}

			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				time.Sleep(500 * time.Millisecond)
				cancel()
			}()

			Download(tt.opt, ctx, client, server.URL+"/file.txt", dir, "file.txt")

			_, err := os.Stat(filepath.Join(dir, "file.txt"))
			exists := err == nil

			if exists != tt.shouldExist {
				t.Fatalf("file existence = %v, expected %v", exists, tt.shouldExist)
			}
		})
	}
}
