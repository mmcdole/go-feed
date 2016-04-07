package atom_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmcdole/gofeed/atom"
	"github.com/stretchr/testify/assert"
)

func TestAtomParser_ParseFeed(t *testing.T) {
	files, _ := filepath.Glob("../testdata/parser/atom/*.xml")
	for _, f := range files {
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		fmt.Printf("Testing %s... ", name)

		// Get actual source feed
		ff := fmt.Sprintf("../testdata/parser/atom/%s.xml", name)
		f, _ := os.Open(ff)
		defer f.Close()

		// Parse actual feed
		fp := &atom.Parser{}
		actual, _ := fp.ParseFeed(f)

		// Get json encoded expected feed result
		ef := fmt.Sprintf("../testdata/parser/atom/%s.json", name)
		e, _ := ioutil.ReadFile(ef)

		// Unmarshal expected feed
		expected := &atom.Feed{}
		json.Unmarshal(e, &expected)

		if assert.Equal(t, actual, expected, "Feed file %s.xml did not match expected output %s.json", name, name) {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("Failed\n")
		}
	}
}
