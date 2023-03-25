package universal

import (
	"sort"
	"testing"
	"time"
)

func TestFeedSort(t *testing.T) {
	oldestItem := &Item{
		PublishedParsed: &[]time.Time{time.Unix(0, 0)}[0],
	}
	inbetweenItem := &Item{
		PublishedParsed: &[]time.Time{time.Unix(1, 0)}[0],
	}
	newestItem := &Item{
		PublishedParsed: &[]time.Time{time.Unix(2, 0)}[0],
	}

	feed := Feed{
		Items: []*Item{
			newestItem,
			oldestItem,
			inbetweenItem,
		},
	}
	expected := Feed{
		Items: []*Item{
			oldestItem,
			inbetweenItem,
			newestItem,
		},
	}

	sort.Sort(feed)

	for i, item := range feed.Items {
		if !item.PublishedParsed.Equal(
			*expected.Items[i].PublishedParsed,
		) {
			t.Errorf(
				"Item PublishedParsed = %s; want %s",
				item.PublishedParsed,
				expected.Items[i].PublishedParsed,
			)
		}
	}
}
