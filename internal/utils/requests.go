package utils

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/google/brotli/go/cbrotli"
)

func Get(ctx context.Context, url string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to create request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to get data: %w", err)
	}
	defer resp.Body.Close()
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to create gzip reader: %w", err)
		}
	case "br":
		reader = cbrotli.NewReader(resp.Body)
	default:
		reader = resp.Body
	}
	defer reader.Close()
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read response body: %w", err)
	}
	return b, resp.StatusCode, nil
}
