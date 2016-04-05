package gofeed

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mmcdole/gofeed/atom"
	"github.com/mmcdole/gofeed/rss"
	"github.com/mmcdole/goxpp"
)

// FeedType represents one of the possible feed
// types that we can detect.
type FeedType int

const (
	// FeedTypeUnknown represents a feed that could not have its
	// type determiend.
	FeedTypeUnknown FeedType = iota
	// FeedTypeAtom repesents an Atom feed
	FeedTypeAtom
	// FeedTypeRSS represents an RSS feed
	FeedTypeRSS
)

// DetectFeedType takes a feed XML string and attempts
// to detect its feed type.
func DetectFeedType(feed string) FeedType {
	p := xpp.NewXMLPullParser(strings.NewReader(feed), false)

	_, err := p.NextTag()
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

// FeedParser is a universal feed parser that detects
// a given feed type, parsers it, and translates it
// to the universal feed type.
type FeedParser struct {
	AtomTrans AtomTranslator
	RSSTrans  RSSTranslator
	rp        *rss.Parser
	ap        *atom.Parser
}

// NewFeedParser creates a FeedParser.
func NewFeedParser() *FeedParser {
	fp := FeedParser{
		rp: &rss.Parser{},
		ap: &atom.Parser{},
	}
	return &fp
}

// ParseFeedURL fetches the contents of a given feed url and
// parses the feed into the universal feed type.
func (f *FeedParser) ParseFeedURL(feedURL string) (*Feed, error) {
	resp, err := http.Get(feedURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return f.ParseFeed(string(body)), nil
}

// ParseFeed takes a feed XML string and parses it into the
// universal feed type.
func (f *FeedParser) ParseFeed(feed string) (*Feed, error) {
	fmt.Println(feed)
	ft := DetectFeedType(feed)
	switch ft {
	case FeedTypeAtom:
		return f.parseFeedFromAtom(feed)
	case FeedTypeRSS:
		return f.parseFeedFromRSS(feed)
	}
	return nil, errors.New("Failed to detect feed type")
}

func (f *FeedParser) parseFeedFromAtom(feed string) (*Feed, error) {
	af, err := f.ap.ParseFeed(feed)
	if err != nil {
		return nil, err
	}
	result := f.atomTrans().Translate(af)
	return result, nil
}

func (f *FeedParser) parseFeedFromRSS(feed string) (*Feed, error) {
	rf, err := f.rp.ParseFeed(feed)
	if err != nil {
		return nil, err
	}

	result := f.rssTrans().Translate(rf)
	return result, nil
}

func (f *FeedParser) atomTrans() AtomTranslator {
	if f.AtomTrans != nil {
		return f.AtomTrans
	}
	f.AtomTrans = &DefaultAtomTranslator{}
	return f.AtomTrans
}

func (f *FeedParser) rssTrans() RSSTranslator {
	if f.RSSTrans != nil {
		return f.RSSTrans
	}
	f.RSSTrans = &DefaultRSSTranslator{}
	return f.RSSTrans
}
