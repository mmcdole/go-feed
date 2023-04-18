package universal

import (
	jsonEncoding "encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmcdole/gofeed/v2/parsers/atom"
	"github.com/stretchr/testify/assert"
)

func TestDefaultAtomConverter_Convert(t *testing.T) {
	files, _ := filepath.Glob("testdata/parser/translator/atom/*.xml")
	for _, f := range files {
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		fmt.Printf("Testing %s... ", name)

		// Get actual source feed
		ff := fmt.Sprintf("testdata/parser/translator/atom/%s.xml", name)
		f, _ := os.Open(ff)
		defer f.Close()

		// Parse actual feed
		translator := &DefaultAtomConverter{}
		fp := &atom.Parser{}
		atomFeed, _ := fp.Parse(f)
		actual, _ := translator.Convert(atomFeed)

		// Get json encoded expected feed result
		ef := fmt.Sprintf("testdata/parser/translator/atom/%s.json", name)
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

func TestDefaultAtomConverter_Convert_WrongType(t *testing.T) {
	translator := &DefaultAtomConverter{}
	af, err := translator.Convert("wrong type")
	assert.Nil(t, af)
	assert.NotNil(t, err)
}
