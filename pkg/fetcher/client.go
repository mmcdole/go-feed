package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Config struct {
	Timeout   time.Duration
	UserAgent string
}

// Response is a struct that holds the response from a feed fetch
type Response struct {
	Body io.ReadCloser
	Err  error
}

// FetchURL retrieves the feed at a given URL using the provided config
func FetchURL(ctx context.Context, url string, config *Config) *Response {
	// Set default values for config if not provided
	if config == nil {
		config = &Config{
			Timeout:   10 * time.Second,
			UserAgent: "gofeed",
		}
	}

	client := &http.Client{
		Timeout: config.Timeout,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return &Response{Err: err}
	}

	req.Header.Set("User-Agent", config.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return &Response{Err: err}
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return &Response{Err: fmt.Errorf("received non-200 status code: %d", resp.StatusCode)}
	}

	return &Response{Body: resp.Body}
}
