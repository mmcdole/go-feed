package config

import (
	"context"
	"time"
)

type BasicAuth struct {
	Username string
	Password string
}

type RequestOptions struct {
	// Timeout is the maximum amount of time to wait for the request to complete.
	// The default value is 29 seconds.
	Timeout time.Duration

	// UserAgent is the value of the `User-Agent` header to be sent in the request.
	// This header is used to identify the client software and version to the server.
	// The default value is "gofeed/1.0.0" (the current version of gofeed).
	UserAgent string

	// IfNoneMatch is the value of the `If-None-Match` header to be sent in the request.
	// This header is used to make a conditional request for a resource, and the server will return the resource only if it does not match the provided entity tag.
	// If the server has a matching `Etag`, it will not respond with the full feed, and instead return a `303 Not Modified` response.
	// The default value is an empty string, which means no `If-None-Match` header will be sent.
	IfNoneMatch string

	// IfModifiedSince is the value of the `If-Modified-Since` header to be sent in the request.
	// This header is used to make a conditional request for a resource, and the server will return the resource only if it has been modified since the provided time.
	// If the server determines that the resource has not been modified since the provided `If-Modified-Since` time, it will not respond with the full feed, and instead return a `303 Not Modified` response.
	// The default value is an empty string, which means no `If-Modified-Since` header will be sent.
	IfModifiedSince string

	// ProxyURL is the optional URL of the proxy to be used for requests.
	// The default value is an empty string, which means no proxy will be used.
	ProxyURL string

	// ProxyAuth is the optional basic authentication settings for the proxy.
	ProxyAuth *BasicAuth

	// FeedAuth is the optional basic authentication settings for accessing the feed.
	FeedAuth *BasicAuth

	// Context is the context object to be passed to the request handlers.
	Context context.Context
}

func NewRequestOptions() RequestOptions {
	return RequestOptions{
		Timeout:         time.Duration(29) * time.Second,
		UserAgent:       "gofeed/2.0.0",
		IfNoneMatch:     "",
		IfModifiedSince: "",
		Context:         context.Background(),
		ProxyURL:        "",
		ProxyAuth:       nil,
		FeedAuth:        nil,
	}
}
