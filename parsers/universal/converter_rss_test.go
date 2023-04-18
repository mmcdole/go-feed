package universal

import (
	jsonEncoding "encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmcdole/gofeed/v2/parsers/rss"
	"github.com/stretchr/testify/assert"
)

func TestDefaultRSSConverter_Convert(t *testing.T) {
	files, _ := filepath.Glob("testdata/parser/translator/rss/*.xml")
	for _, f := range files {
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		fmt.Printf("Testing %s... ", name)

		// Get actual source feed
		ff := fmt.Sprintf("testdata/parser/translator/rss/%s.xml", name)
		f, _ := os.Open(ff)
		defer f.Close()

		// Parse actual feed
		translator := &DefaultRSSConverter{}
		fp := &rss.Parser{}
		rssFeed, _ := fp.Parse(f)
		actual, _ := translator.Convert(rssFeed)

		// Get json encoded expected feed result
		ef := fmt.Sprintf("testdata/parser/translator/rss/%s.json", name)
		e, _ := ioutil.ReadFile(ef)

		// Unmarshal expected feed
		expected := &Feed{}
		jsonEncoding.Unmarshal(e, &expected)

		if assert.Equal(t, expected, actual, "Feed file %s.xml did not match expected output %s.json", name, name) {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("Failed\n")
		}
	}
}

func TestDefaultRSSConverter_Convert_WrongType(t *testing.T) {
	translator := &DefaultRSSConverter{}
	af, err := translator.Convert("wrong type")
	assert.Nil(t, af)
	assert.NotNil(t, err)
}
