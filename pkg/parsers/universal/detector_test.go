package universal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectFeedType(t *testing.T) {
	var feedTypeTests = []struct {
		file     string
		expected FeedType
	}{
		{"atom03_feed.xml", FeedTypeAtom},
		{"atom10_feed.xml", FeedTypeAtom},
		{"rss_feed.xml", FeedTypeRSS},
		{"rss_feed_bom.xml", FeedTypeRSS},
		{"rss_feed_leading_spaces.xml", FeedTypeRSS},
		{"rdf_feed.xml", FeedTypeRSS},
		{"unknown_feed.xml", FeedTypeUnknown},
		{"empty_feed.xml", FeedTypeUnknown},
		{"json10_feed.json", FeedTypeJSON},
	}

	for _, test := range feedTypeTests {
		fmt.Printf("Testing %s... ", test.file)

		// Get feed content
		path := fmt.Sprintf("../../../testdata/parser/universal/%s", test.file)
		f, _ := ioutil.ReadFile(path)

		// Get actual value
		actual := DetectFeedType(bytes.NewReader(f))

		if assert.Equal(t, actual, test.expected, "Feed file %s did not match expected type %d", test.file, test.expected) {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("Failed\n")
		}
	}
}

// Examples

func ExampleDetectFeedType() {
	feedData := `<rss version="2.0">
<channel>
<title>Sample Feed</title>
</channel>
</rss>`
	feedType := DetectFeedType(strings.NewReader(feedData))
	if feedType == FeedTypeRSS {
		fmt.Println("Wow! This is an RSS feed!")
	}
}
