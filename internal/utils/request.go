package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type CurlConf[T any] struct {
	Decoder func(io.Reader) error
}

type CurlOpts[T any] func(*CurlConf[T])

func WithDecoder[T any](decoder func(io.Reader) error) func(*CurlConf[T]) {
	return func(c *CurlConf[T]) {
		c.Decoder = decoder
	}
}

func PatchCurl[T any](ctx context.Context, targetURL string, payload any, headers map[string]string) (res *T, statusCode int, err error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return res, 0, fmt.Errorf("unable to marshal payload: %w", err)
	}
	return CurlWithBody[T](ctx, http.MethodPatch, targetURL, payloadJSON, headers)
}

func PostCurl[T any](ctx context.Context, targetURL string, payload any, headers map[string]string) (res *T, statusCode int, err error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return res, 0, fmt.Errorf("unable to marshal payload: %w", err)
	}
	return CurlWithBody[T](ctx, http.MethodPost, targetURL, payloadJSON, headers)
}

func CurlWithBody[T any](ctx context.Context, method, targetURL string, payloadJSON []byte, headers map[string]string) (res *T, statusCode int, err error) {
	// todo inject tracer
	req, err := http.NewRequestWithContext(ctx, method, targetURL, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return res, 0, fmt.Errorf("unable to create request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return executeWithDefaultClient[T](req, nil)
}

// GetCurl is a generic function to send GET request with headers and return response
// expect json response
// T is a type of response, automatically unmarshalled from json
// check status code first. in some cases response can be different, so unmarsharlling will fail
// status code return as is even in unmarshalling error
func GetCurl[T any](ctx context.Context, targetURL string, headers map[string]string, opts ...CurlOpts[T]) (res *T, statusCode int, err error) {
	config := &CurlConf[T]{}
	for _, opt := range opts {
		opt(config)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, http.NoBody)
	if err != nil {
		return res, 0, fmt.Errorf("unable to create request: %w", err)
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return executeWithDefaultClient[T](req, config.Decoder)
}

func executeWithDefaultClient[T any](req *http.Request, decoder func(io.Reader) error) (res *T, statusCode int, err error) {
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return res, 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if decoder != nil {
		return nil, resp.StatusCode, decoder(resp.Body)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, 0, fmt.Errorf("failed to read response body: %w", err)
	}
	var result T
	if strings.TrimSpace(strings.ReplaceAll(string(b), "\"", "")) != "" {
		if err = json.Unmarshal(b, &result); err != nil {
			return res, resp.StatusCode, fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return &result, resp.StatusCode, nil
}
