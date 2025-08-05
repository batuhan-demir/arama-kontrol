package file

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func Download(url string, output string) error {
	dir := "./files"

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Add user agent to avoid blocks
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download: %v", err)
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	// Check Content-Length
	if resp.ContentLength == 0 {
		return fmt.Errorf("file is empty (Content-Length: 0)")
	}

	filePath := fmt.Sprintf("%s/%s", dir, output)

	// Create temp file first
	tempFile := filePath + ".tmp"
	out, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}

	// Copy with progress tracking
	written, err := io.Copy(out, resp.Body)
	out.Close()

	if err != nil {
		os.Remove(tempFile) // Clean up temp file
		return fmt.Errorf("failed to copy data: %v", err)
	}

	// Check if we actually wrote data
	if written == 0 {
		os.Remove(tempFile)
		return fmt.Errorf("no data written to file")
	}

	// Move temp file to final location
	err = os.Rename(tempFile, filePath)
	if err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to move temp file: %v", err)
	}

	fmt.Printf("Downloaded %s: %d bytes\n", output, written)
	return nil
}
