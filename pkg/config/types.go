package config

import "time"

type ParseOptions struct {
	// MaxItems specifies the maximum number of feed items to parse.
	// The default value is 0, which means no limit.
	MaxItems int

	// ParseDates determines if the feed parser will attempt to parse dates into `time.Time` objects.
	// The default value is true.
	ParseDates bool

	// ParseExtensions determines if the feed parser will attempt to parse feed extensions such as
	// iTunes extensions, or custom feed extensions.
	// The default value is true.
	ParseExtensions bool

	// KeepOriginalFeed specifies if the parser should retain the raw feed in the `Feed` struct's `RawFeed` field.
	// The default value is false.
	KeepOriginalFeed bool

	// StrictMode determines if strict parsing rules should be applied.
	// If set to true, all strictness settings are enabled.
	// The default value is false.
	StrictMode bool

	// StrictnessOptions holds the options for controlling the strictness of the parsing.
	StrictnessOptions StrictnessOptions

	// RequestOptions holds the options for the HTTP request.
	RequestOptions RequestOptions
}

type StrictnessOptions struct {
	// StripInvalidCharacters specifies if invalid feed characters should be stripped out.
	// The default value is false.
	StripInvalidCharacters bool

	// AutoCloseTags specifies if the parser should automatically close unclosed tags.
	// The default value is false.
	AutoCloseTags bool

	// AllowUndisclosedXMLNamespaces specifies if the parser should allow undisclosed XML namespaces.
	// The default value is false.
	AllowUndisclosedXMLNamespaces bool

	// AllowIncorrectDateFormats specifies if the parser should allow incorrect date formats.
	// The default value is false.
	AllowIncorrectDateFormats bool

	// AllowUnescapedMarkup specifies if the parser should allow unescaped / naked markup in feed elements.
	// The default value is false.
	AllowUnescapedMarkup bool
}

type RequestOptions struct {
	// Timeout is the maximum amount of time to wait for the request to complete.
	// The default value is 0, which means no timeout.
	Timeout time.Duration

	// UserAgent is the value of the `User-Agent` header to be sent in the request.
	// This header is used to identify the client software and version to the server.
	// The default value is "gofeed/2.0.0" (the current version of gofeed).
	UserAgent string

	// IfNoneMatch is the value of the `If-None-Match` header to be sent in the request.
	// This header is used to make a conditional request for a resource, and the server will return the resource only if it does not match the provided entity tag.
	// If the server has a matching `Etag`, it will not respond with the full feed, and instead return a `304 Not Modified` response.
	// The default value is an empty string, which means no `If-None-Match` header will be sent.
	IfNoneMatch string

	// IfModifiedSince is the value of the `If-Modified-Since` header to be sent in the request.
	// This header is used to make a conditional request for a resource, and the server will return the resource only if it has been modified since the provided time.
	// If the server determines that the resource has not been modified since the provided `If-Modified-Since` time, it will not respond with the full feed, and instead return a `304 Not Modified` response.
	// The default value is an empty string, which means no `If-Modified-Since` header will be sent.
	IfModifiedSince string
}
