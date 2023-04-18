package universal

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"github.com/mmcdole/gofeed/v2/config"
	"github.com/mmcdole/gofeed/v2/internal/shared"
	"github.com/mmcdole/gofeed/v2/parsers/atom"
	"github.com/mmcdole/gofeed/v2/parsers/json"
	"github.com/mmcdole/gofeed/v2/parsers/rss"
)

type Converter interface {
	Convert(feed interface{}) (*Feed, error)
}

// ErrFeedTypeNotDetected is returned when the detection system can not figure
// out the Feed format
var ErrFeedTypeNotDetected = errors.New("failed to detect feed type")

// Parser is a universal feed parser that can parse Atom, RSS, and JSON feeds.
type Parser struct {
	atomParser atom.Parser
	rssParser  rss.Parser
	jsonParser json.Parser

	AtomConverter Converter
	RSSConverter  Converter
	JSONConverter Converter
}

func NewParser() *Parser {
	return &Parser{
		atomParser: *atom.NewParser(),
		rssParser:  *rss.NewParser(),
		jsonParser: *json.NewParser(),
	}
}

// Parse reads an XML, JSON, RSS, or Atom feed from an io.Reader and parses it
// into a universal Feed struct using default options.
func (f *Parser) Parse(feed io.Reader) (*Feed, error) {
	return f.ParseWithOptions(feed, config.NewParseOptions())
}

func (f *Parser) ParseWithOptions(feed io.Reader, options config.ParseOptions) (*Feed, error) {
	// Wrap the feed io.Reader in an io.TeeReader to capture all the bytes
	// read by the DetectFeedType function and construct a new reader with
	// those bytes intact for when we attempt to parse the feeds.
	var buf bytes.Buffer
	tee := io.TeeReader(feed, &buf)
	feedType := DetectFeedType(tee)

	// Glue the read bytes from the detect function back into a new reader
	reader := io.MultiReader(&buf, feed)

	switch feedType {
	case FeedTypeAtom:
		parsedFeed, err := f.atomParser.ParseWithOptions(reader, options)
		if err != nil {
			return nil, err
		}
		return f.getAtomConverter().Convert(parsedFeed)
	case FeedTypeRSS:
		parsedFeed, err := f.rssParser.ParseWithOptions(reader, options)
		if err != nil {
			return nil, err
		}
		return f.getRSSConverter().Convert(parsedFeed)
	case FeedTypeJSON:
		parsedFeed, err := f.jsonParser.ParseWithOptions(reader, options)
		if err != nil {
			return nil, err
		}
		return f.getJSONConverter().Convert(parsedFeed)
	default:
		return nil, ErrFeedTypeNotDetected
	}
}

// ParseURL fetches the contents of a given URL and attempts to parse the
// response into a universal Feed struct using default RequestOptions and
// ParseOptions.
func (f *Parser) ParseURL(feedURL string) (*Feed, error) {
	return f.ParseURLWithOptions(feedURL, config.NewRequestOptions(), config.NewParseOptions())
}

// ParseURLWithOptions fetches the contents of a given URL and attempts to
// parse the response into a universal Feed struct. It takes a string URL,
// RequestOptions for custom request behavior, and ParseOptions for custom
// parsing behavior.
func (f *Parser) ParseURLWithOptions(feedURL string, reqOptions config.RequestOptions, parseOptions config.ParseOptions) (feed *Feed, err error) {
	reqClient := &shared.RequestClient{}
	body, err := reqClient.Fetch(feedURL, reqOptions)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return f.ParseWithOptions(body, parseOptions)
}

// ParseString parses a feed string into the universal feed type using default
// ParseOptions.
func (f *Parser) ParseString(feed string) (*Feed, error) {
	return f.ParseStringWithOptions(feed, config.NewParseOptions())
}

// ParseStringWithOptions parses a feed string into the universal feed type
// using the provided ParseOptions.
func (f *Parser) ParseStringWithOptions(feed string, options config.ParseOptions) (*Feed, error) {
	return f.Parse(strings.NewReader(feed))
}

func (f *Parser) getAtomConverter() Converter {
	if f.AtomConverter == nil {
		f.AtomConverter = &DefaultAtomConverter{}
	}
	return f.AtomConverter
}

func (f *Parser) getRSSConverter() Converter {
	if f.RSSConverter == nil {
		f.RSSConverter = &DefaultRSSConverter{}
	}
	return f.RSSConverter
}

func (f *Parser) getJSONConverter() Converter {
	if f.JSONConverter == nil {
		f.JSONConverter = &DefaultJSONConverter{}
	}
	return f.JSONConverter
}
