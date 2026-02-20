package downloader

import (
	"net/http"
	"net/http/httptest"
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

	err := Download(client, server.URL+"/file.txt", dir, "fallback.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
