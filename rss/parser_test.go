package rss_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmcdole/gofeed/rss"
	"github.com/stretchr/testify/assert"
)

func TestRSSParser_ParseFeed(t *testing.T) {
	files, _ := filepath.Glob("../testdata/parser/rss/*.xml")
	for _, f := range files {
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		fmt.Printf("Testing %s... ", name)

		// Get actual source feed
		ff := fmt.Sprintf("../testdata/parser/rss/%s.xml", name)
		f, _ := os.Open(ff)
		defer f.Close()

		// Parse actual feed
		fp := &rss.Parser{}
		actual, _ := fp.ParseFeed(f)

		// Get json encoded expected feed result
		ef := fmt.Sprintf("../testdata/parser/rss/%s.json", name)
		e, _ := ioutil.ReadFile(ef)

		// Unmarshal expected feed
		expected := &rss.Feed{}
		json.Unmarshal(e, &expected)

		if assert.Equal(t, expected, actual, "Feed file %s.xml did not match expected output %s.json", name, name) {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("Failed\n")
		}
	}
}
