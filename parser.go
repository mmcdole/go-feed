package gofeed

import (
	"github.com/mmcdole/gofeed/v2/parsers/universal"
)

// GofeedParser is a convienance wrapper around the universal parser.
type GofeedParser struct {
	*universal.Parser
}

// NewParser returns a new GofeedParser.
func NewParser() *GofeedParser {
	return &GofeedParser{
		Parser: universal.NewParser(),
	}
}
