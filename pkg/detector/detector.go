package gofeed

import (
	"bufio"
	"io"
	"strings"
	"unicode"

	jsoniter "github.com/json-iterator/go"
	xpp "github.com/mmcdole/goxpp"

	"github.com/mmcdole/gofeed/v2/internal/shared"
)

const (
	// FeedTypeUnknown represents a feed that could not have its
	// type determined.
	FeedTypeUnknown FeedType = iota
	// FeedTypeAtom represents an Atom feed
	FeedTypeAtom
	// FeedTypeRSS represents an RSS feed
	FeedTypeRSS
	// FeedTypeJSON represents a JSON feed
	FeedTypeJSON
)

const (
	jsonStartChar = '{'
	xmlStartChar  = '<'
	bufSize       = 1024
)

// DetectFeedType attempts to determine the type of feed
// by looking for specific xml elements unique to the
// various feed types.
func DetectFeedType(feed io.Reader) FeedType {
	bufReader := bufio.NewReaderSize(feed, bufSize)

	var firstChar byte
	for {
		ch, _, err := bufReader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return FeedTypeUnknown
		}
		// ignore leading whitespace & byte order marks
		if isWhitespaceOrBOM(ch) {
			continue
		}
		firstChar = byte(ch)
		break
	}

	switch firstChar {
	case jsonStartChar:
		// Check if document is valid JSON
		n := bufReader.Buffered()
		b := make([]byte, n)
		_, err := bufReader.Read(b)
		if err != nil {
			return FeedTypeUnknown
		}
		if jsoniter.Valid(b) {
			return FeedTypeJSON
		}
	case xmlStartChar:
		return detectXMLFeedType(bufReader)
	}
	return FeedTypeUnknown
}

func detectXMLFeedType(bufReader *bufio.Reader) FeedType {
	p := xpp.NewXMLPullParser(bufReader, false, shared.NewReaderLabel)

	xmlBase := shared.XMLBase{}
	_, err := xmlBase.FindRoot(p)
	if err != nil {
		return FeedTypeUnknown
	}

	name := strings.ToLower(p.Name)
	switch name {
	case "rdf":
		return FeedTypeRSS
	case "rss":
		return FeedTypeRSS
	case "feed":
		return FeedTypeAtom
	default:
		return FeedTypeUnknown
	}
}

func isWhitespaceOrBOM(ch rune) bool {
	return unicode.IsSpace(ch) || isBOM(ch)
}

func isBOM(ch rune) bool {
	switch ch {
	case 0xFE, 0xFF, 0x00, 0xEF, 0xBB, 0xBF: // utf 8-16-32 BOM
		return true
	default:
		return false
	}
}
