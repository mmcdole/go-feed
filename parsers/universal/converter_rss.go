package universal

import (
	"strings"
	"time"

	"github.com/mmcdole/gofeed/v2/config"
	"github.com/mmcdole/gofeed/v2/internal/shared"
	"github.com/mmcdole/gofeed/v2/parsers/universal"
	ext "github.com/mmcdole/gofeed/v2/pkg/extensions"
)

// DefaultRSSConverter converts an rss.Feed struct
// into the generic Feed struct.
//
// This default implementation defines a set of
// mapping rules between rss.Feed -> Feed
// for each of the fields in Feed.
type DefaultRSSConverter struct{}

// convert converts an RSS feed into the universal
// feed type.
func (t *DefaultRSSConverter) convert(feed *Feed, options config.ParseOptions) (*universal.Feed, error) {
	result := &universal.Feed{}
	result.Title = t.convertFeedTitle(feed)
	result.Description = t.convertFeedDescription(feed)
	result.Link = t.convertFeedLink(feed)
	result.Links = t.convertFeedLinks(feed)
	result.FeedLink = t.convertFeedFeedLink(feed)
	result.Updated = t.convertFeedUpdated(feed)
	result.UpdatedParsed = t.convertFeedUpdatedParsed(feed)
	result.Published = t.convertFeedPublished(feed)
	result.PublishedParsed = t.convertFeedPublishedParsed(feed)
	result.Authors = t.convertFeedAuthors(feed)
	result.Language = t.convertFeedLanguage(feed)
	result.Image = t.convertFeedImage(feed)
	result.Copyright = t.convertFeedCopyright(feed)
	result.Generator = t.convertFeedGenerator(feed)
	result.Categories = t.convertFeedCategories(feed)
	result.Items = t.convertFeedItems(feed)
	result.ITunesExt = feed.ITunesExt
	result.DublinCoreExt = feed.DublinCoreExt
	result.Extensions = feed.Extensions
	result.FeedVersion = feed.Version
	result.FeedType = "rss"
	return result, nil
}

func (t *DefaultRSSConverter) convertFeedItem(rssItem *Item) (item *universal.Item) {
	item = &universal.Item{}
	item.Title = t.convertItemTitle(rssItem)
	item.Description = t.convertItemDescription(rssItem)
	item.Content = t.convertItemContent(rssItem)
	item.Link = t.convertItemLink(rssItem)
	item.Links = t.convertItemLinks(rssItem)
	item.Published = t.convertItemPublished(rssItem)
	item.PublishedParsed = t.convertItemPublishedParsed(rssItem)
	item.Authors = t.convertItemAuthors(rssItem)
	item.GUID = t.convertItemGUID(rssItem)
	item.Image = t.convertItemImage(rssItem)
	item.Categories = t.convertItemCategories(rssItem)
	item.Enclosures = t.convertItemEnclosures(rssItem)
	item.DublinCoreExt = rssItem.DublinCoreExt
	item.ITunesExt = rssItem.ITunesExt
	item.Extensions = rssItem.Extensions
	return
}

func (t *DefaultRSSConverter) convertFeedTitle(rss *Feed) (title string) {
	if rss.Title != "" {
		title = rss.Title
	} else if rss.DublinCoreExt != nil && rss.DublinCoreExt.Title != nil {
		title = t.firstEntry(rss.DublinCoreExt.Title)
	}
	return
}

func (t *DefaultRSSConverter) convertFeedDescription(rss *Feed) (desc string) {
	return rss.Description
}

func (t *DefaultRSSConverter) convertFeedLink(rss *Feed) (link string) {
	if rss.Link != "" {
		link = rss.Link
	} else if rss.ITunesExt != nil && rss.ITunesExt.Subtitle != "" {
		link = rss.ITunesExt.Subtitle
	}
	return
}

func (t *DefaultRSSConverter) convertFeedFeedLink(rss *Feed) (link string) {
	atomExtensions := t.extensionsForKeys([]string{"atom", "atom10", "atom03"}, rss.Extensions)
	for _, ex := range atomExtensions {
		if links, ok := ex["link"]; ok {
			for _, l := range links {
				if l.Attrs["rel"] == "self" {
					link = l.Attrs["href"]
				}
			}
		}
	}
	return
}

func (t *DefaultRSSConverter) convertFeedLinks(rss *Feed) (links []string) {
	if len(rss.Links) > 0 {
		links = append(links, rss.Links...)
	}
	atomExtensions := t.extensionsForKeys([]string{"atom", "atom10", "atom03"}, rss.Extensions)
	for _, ex := range atomExtensions {
		if lks, ok := ex["link"]; ok {
			for _, l := range lks {
				if l.Attrs["rel"] == "" || l.Attrs["rel"] == "alternate" || l.Attrs["rel"] == "self" {
					links = append(links, l.Attrs["href"])
				}
			}
		}
	}
	return
}

func (t *DefaultRSSConverter) convertFeedUpdated(rss *Feed) (updated string) {
	if rss.LastBuildDate != "" {
		updated = rss.LastBuildDate
	} else if rss.DublinCoreExt != nil && rss.DublinCoreExt.Date != nil {
		updated = t.firstEntry(rss.DublinCoreExt.Date)
	}
	return
}

func (t *DefaultRSSConverter) convertFeedUpdatedParsed(rss *Feed) (updated *time.Time) {
	if rss.LastBuildDateParsed != nil {
		updated = rss.LastBuildDateParsed
	} else if rss.DublinCoreExt != nil && rss.DublinCoreExt.Date != nil {
		dateText := t.firstEntry(rss.DublinCoreExt.Date)
		date, err := shared.ParseDate(dateText)
		if err == nil {
			updated = &date
		}
	}
	return
}

func (t *DefaultRSSConverter) convertFeedPublished(rss *Feed) (published string) {
	return rss.PubDate
}

func (t *DefaultRSSConverter) convertFeedPublishedParsed(rss *Feed) (published *time.Time) {
	return rss.PubDateParsed
}

func (t *DefaultRSSConverter) convertFeedAuthor(rss *Feed) (author *Person) {
	if rss.ManagingEditor != "" {
		name, address := shared.ParseNameAddress(rss.ManagingEditor)
		author = &Person{}
		author.Name = name
		author.Email = address
	} else if rss.WebMaster != "" {
		name, address := shared.ParseNameAddress(rss.WebMaster)
		author = &Person{}
		author.Name = name
		author.Email = address
	} else if rss.DublinCoreExt != nil && rss.DublinCoreExt.Author != nil {
		dcAuthor := t.firstEntry(rss.DublinCoreExt.Author)
		name, address := shared.ParseNameAddress(dcAuthor)
		author = &Person{}
		author.Name = name
		author.Email = address
	} else if rss.DublinCoreExt != nil && rss.DublinCoreExt.Creator != nil {
		dcCreator := t.firstEntry(rss.DublinCoreExt.Creator)
		name, address := shared.ParseNameAddress(dcCreator)
		author = &Person{}
		author.Name = name
		author.Email = address
	} else if rss.ITunesExt != nil && rss.ITunesExt.Author != "" {
		name, address := shared.ParseNameAddress(rss.ITunesExt.Author)
		author = &Person{}
		author.Name = name
		author.Email = address
	}
	return
}

func (t *DefaultRSSConverter) convertFeedAuthors(rss *Feed) (authors []*Person) {
	if author := t.convertFeedAuthor(rss); author != nil {
		authors = []*Person{author}
	}
	return
}

func (t *DefaultRSSConverter) convertFeedLanguage(rss *Feed) (language string) {
	if rss.Language != "" {
		language = rss.Language
	} else if rss.DublinCoreExt != nil && rss.DublinCoreExt.Language != nil {
		language = t.firstEntry(rss.DublinCoreExt.Language)
	}
	return
}

func (t *DefaultRSSConverter) convertFeedImage(rss *Feed) (image *Image) {
	if rss.Image != nil {
		image = &Image{}
		image.Title = rss.Image.Title
		image.URL = rss.Image.URL
	} else if rss.ITunesExt != nil && rss.ITunesExt.Image != "" {
		image = &Image{}
		image.URL = rss.ITunesExt.Image
	}
	return
}

func (t *DefaultRSSConverter) convertFeedCopyright(rss *Feed) (rights string) {
	if rss.Copyright != "" {
		rights = rss.Copyright
	} else if rss.DublinCoreExt != nil && rss.DublinCoreExt.Rights != nil {
		rights = t.firstEntry(rss.DublinCoreExt.Rights)
	}
	return
}

func (t *DefaultRSSConverter) convertFeedGenerator(rss *Feed) (generator string) {
	return rss.Generator
}

func (t *DefaultRSSConverter) convertFeedCategories(rss *Feed) (categories []string) {
	cats := []string{}
	if rss.Categories != nil {
		for _, c := range rss.Categories {
			cats = append(cats, c.Value)
		}
	}

	if rss.ITunesExt != nil && rss.ITunesExt.Keywords != "" {
		keywords := strings.Split(rss.ITunesExt.Keywords, ",")
		for _, k := range keywords {
			cats = append(cats, k)
		}
	}

	if rss.ITunesExt != nil && rss.ITunesExt.Categories != nil {
		for _, c := range rss.ITunesExt.Categories {
			cats = append(cats, c.Text)
			if c.Subcategory != nil {
				cats = append(cats, c.Subcategory.Text)
			}
		}
	}

	if rss.DublinCoreExt != nil && rss.DublinCoreExt.Subject != nil {
		for _, c := range rss.DublinCoreExt.Subject {
			cats = append(cats, c)
		}
	}

	if len(cats) > 0 {
		categories = cats
	}

	return
}

func (t *DefaultRSSConverter) convertFeedItems(rss *Feed) (items []*Item) {
	items = []*Item{}
	for _, i := range rss.Items {
		items = append(items, t.convertFeedItem(i))
	}
	return
}

func (t *DefaultRSSConverter) convertItemTitle(rssItem *Item) (title string) {
	if rssItem.Title != "" {
		title = rssItem.Title
	} else if rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Title != nil {
		title = t.firstEntry(rssItem.DublinCoreExt.Title)
	}
	return
}

func (t *DefaultRSSConverter) convertItemDescription(rssItem *Item) (desc string) {
	if rssItem.Description != "" {
		desc = rssItem.Description
	} else if rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Description != nil {
		desc = t.firstEntry(rssItem.DublinCoreExt.Description)
	}
	return
}

func (t *DefaultRSSConverter) convertItemContent(rssItem *Item) (content string) {
	return rssItem.Content
}

func (t *DefaultRSSConverter) convertItemLink(rssItem *Item) (link string) {
	return rssItem.Link
}

func (t *DefaultRSSConverter) convertItemLinks(rssItem *Item) (links []string) {
	if len(rssItem.Links) > 0 {
		links = append(links, rssItem.Links...)
	}
	return links
}

func (t *DefaultRSSConverter) convertItemUpdated(rssItem *Item) (updated string) {
	if rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Date != nil {
		updated = t.firstEntry(rssItem.DublinCoreExt.Date)
	}
	return updated
}

func (t *DefaultRSSConverter) convertItemUpdatedParsed(rssItem *Item) (updated *time.Time) {
	if rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Date != nil {
		updatedText := t.firstEntry(rssItem.DublinCoreExt.Date)
		updatedDate, err := shared.ParseDate(updatedText)
		if err == nil {
			updated = &updatedDate
		}
	}
	return
}

func (t *DefaultRSSConverter) convertItemPublished(rssItem *Item) (pubDate string) {
	if rssItem.PubDate != "" {
		return rssItem.PubDate
	} else if rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Date != nil {
		return t.firstEntry(rssItem.DublinCoreExt.Date)
	}
	return
}

func (t *DefaultRSSConverter) convertItemPublishedParsed(rssItem *Item) (pubDate *time.Time) {
	if rssItem.PubDateParsed != nil {
		return rssItem.PubDateParsed
	} else if rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Date != nil {
		pubDateText := t.firstEntry(rssItem.DublinCoreExt.Date)
		pubDateParsed, err := shared.ParseDate(pubDateText)
		if err == nil {
			pubDate = &pubDateParsed
		}
	}
	return
}

func (t *DefaultRSSConverter) convertItemAuthor(rssItem *Item) (author *Person) {
	if rssItem.Author != "" {
		name, address := shared.ParseNameAddress(rssItem.Author)
		author = &Person{}
		author.Name = name
		author.Email = address
	} else if rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Author != nil {
		dcAuthor := t.firstEntry(rssItem.DublinCoreExt.Author)
		name, address := shared.ParseNameAddress(dcAuthor)
		author = &Person{}
		author.Name = name
		author.Email = address
	} else if rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Creator != nil {
		dcCreator := t.firstEntry(rssItem.DublinCoreExt.Creator)
		name, address := shared.ParseNameAddress(dcCreator)
		author = &Person{}
		author.Name = name
		author.Email = address
	} else if rssItem.ITunesExt != nil && rssItem.ITunesExt.Author != "" {
		name, address := shared.ParseNameAddress(rssItem.ITunesExt.Author)
		author = &Person{}
		author.Name = name
		author.Email = address
	}
	return
}

func (t *DefaultRSSConverter) convertItemAuthors(rssItem *Item) (authors []*Person) {
	if author := t.convertItemAuthor(rssItem); author != nil {
		authors = []*Person{author}
	}
	return
}

func (t *DefaultRSSConverter) convertItemGUID(rssItem *Item) (guid string) {
	if rssItem.GUID != nil {
		guid = rssItem.GUID.Value
	}
	return
}

func (t *DefaultRSSConverter) convertItemImage(rssItem *Item) (image *Image) {
	if rssItem.ITunesExt != nil && rssItem.ITunesExt.Image != "" {
		image = &Image{}
		image.URL = rssItem.ITunesExt.Image
	}
	return
}

func (t *DefaultRSSConverter) convertItemCategories(rssItem *Item) (categories []string) {
	cats := []string{}
	if rssItem.Categories != nil {
		for _, c := range rssItem.Categories {
			cats = append(cats, c.Value)
		}
	}

	if rssItem.ITunesExt != nil && rssItem.ITunesExt.Keywords != "" {
		keywords := strings.Split(rssItem.ITunesExt.Keywords, ",")
		for _, k := range keywords {
			cats = append(cats, k)
		}
	}

	if rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Subject != nil {
		for _, c := range rssItem.DublinCoreExt.Subject {
			cats = append(cats, c)
		}
	}

	if len(cats) > 0 {
		categories = cats
	}

	return
}

func (t *DefaultRSSConverter) convertItemEnclosures(rssItem *Item) (enclosures []*Enclosure) {
	if rssItem.Enclosures != nil && len(rssItem.Enclosures) > 0 {
		// Accumulate the enclosures
		for _, enc := range rssItem.Enclosures {
			e := &Enclosure{}
			e.URL = enc.URL
			e.Type = enc.Type
			e.Length = enc.Length
			enclosures = append(enclosures, e)
		}
	}

	if len(enclosures) == 0 {
		enclosures = nil
	}

	return
}

func (t *DefaultRSSConverter) extensionsForKeys(keys []string, extensions ext.Extensions) (matches []map[string][]ext.Extension) {
	matches = []map[string][]ext.Extension{}

	if extensions == nil {
		return
	}

	for _, key := range keys {
		if match, ok := extensions[key]; ok {
			matches = append(matches, match)
		}
	}
	return
}

func (t *DefaultRSSConverter) firstEntry(entries []string) (value string) {
	if entries == nil {
		return
	}

	if len(entries) == 0 {
		return
	}

	return entries[0]
}
