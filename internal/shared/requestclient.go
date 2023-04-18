package shared

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/mmcdole/gofeed/v2/config"
)

// RequestClient is a wrapper around an http.Client
type RequestClient struct {
	*http.Client
}

func (rc *RequestClient) Fetch(feedURL string, options config.RequestOptions) (io.ReadCloser, error) {
	rc.Client = &http.Client{
		Timeout: options.Timeout,
	}

	if options.ProxyURL != "" {
		proxyURL, err := url.Parse(options.ProxyURL)
		if err != nil {
			return nil, err
		}
		rc.Client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	req, err := http.NewRequestWithContext(options.Context, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", options.UserAgent)

	if options.IfNoneMatch != "" {
		req.Header.Set("If-None-Match", options.IfNoneMatch)
	}

	if options.IfModifiedSince != "" {
		req.Header.Set("If-Modified-Since", options.IfModifiedSince)
	}

	if options.FeedAuth != nil && options.FeedAuth.Username != "" && options.FeedAuth.Password != "" {
		req.SetBasicAuth(options.FeedAuth.Username, options.FeedAuth.Password)
	}

	if options.ProxyAuth != nil && options.ProxyAuth.Username != "" && options.ProxyAuth.Password != "" {
		basicAuth := options.ProxyAuth.Username + ":" + options.ProxyAuth.Password
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(basicAuth))
		req.Header.Set("Proxy-Authorization", "Basic "+encodedAuth)
	}

	resp, err := rc.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
		}
	}

	return resp.Body, nil
}

// ErrFeedTypeNotDetected is returned when the detection system can not figure
// out the Feed format
var ErrFeedTypeNotDetected = errors.New("failed to detect feed type")

// HTTPError represents an HTTP error returned by a server.
type HTTPError struct {
	StatusCode int
	Status     string
}

func (err HTTPError) Error() string {
	return fmt.Sprintf("http error: %s", err.Status)
}
