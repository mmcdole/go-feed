package json

import (
	"bytes"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/mmcdole/gofeed/v2/config"
	"github.com/mmcdole/gofeed/v2/internal/shared"
)

var (
	j = jsoniter.ConfigCompatibleWithStandardLibrary
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (jp *Parser) Parse(feed io.Reader) (*Feed, error) {
	return jp.ParseWithOptions(feed, config.NewParseOptions())
}

func (jp *Parser) ParseWithOptions(feed io.Reader, options config.ParseOptions) (*Feed, error) {
	jsonFeed := &Feed{}

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(feed)

	err := j.Unmarshal(buffer.Bytes(), jsonFeed)
	if err != nil {
		return nil, err
	}
	return jsonFeed, err
}

func (jp *Parser) ParseString(feed string) (*Feed, error) {
	return jp.ParseStringWithOptions(feed, config.NewParseOptions())
}

// ParseString parses an xml feed into an interface{} of atom.Feed
func (jp *Parser) ParseStringWithOptions(feed string, options config.ParseOptions) (*Feed, error) {
	return jp.ParseWithOptions(strings.NewReader(feed), options)
}

func (jp *Parser) ParseURL(feedURL string) (*Feed, error) {
	return jp.ParseURLWithOptions(feedURL, config.NewParseOptions(), config.NewRequestOptions())
}

func (jp *Parser) ParseURLWithOptions(feedURL string, parseOptions config.ParseOptions, reqOptions config.RequestOptions) (*Feed, error) {
	reqClient := &shared.RequestClient{}
	body, err := reqClient.Fetch(feedURL, reqOptions)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return jp.ParseWithOptions(body, parseOptions)
}
