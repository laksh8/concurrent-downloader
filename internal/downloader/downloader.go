package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/laksh8/concurrent-downloader/internal/util"
)

func Download(client *http.Client, rawUrl string, destDir string, fallback string) error {
	req, err := http.NewRequest(http.MethodGet, rawUrl, nil)
	if err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("something went wrong at %s: %s", rawUrl, resp.Status)
	}

	filename := util.SafeFilename(rawUrl, fallback)
	fullPath := filepath.Join(destDir, filename)

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("unable to download to file: %w", err)
	}
	return nil
}
