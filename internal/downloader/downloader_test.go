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

	err := Download(ctx, client, server.URL+"/file.txt", dir, "file.txt")
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

	err := Download(ctx, client, "://bad-url", dir, "fallback.txt")
	if err == nil {
		t.Fatalf("expected error for invalid URL")
	}
}
