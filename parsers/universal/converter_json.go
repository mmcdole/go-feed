package universal

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed/v2/internal/shared"
	"github.com/mmcdole/gofeed/v2/parsers/json"
)

// DefaultJSONConverter converts an json.Feed struct
// into the generic Feed struct.
//
// This default implementation defines a set of
// mapping rules between json.Feed -> Feed
// for each of the fields in Feed.
type DefaultJSONConverter struct{}

// Convert converts an JSON feed into the universal
// feed type.
func (t *DefaultJSONConverter) Convert(feed interface{}) (*Feed, error) {
	json, found := feed.(*json.Feed)
	if !found {
		return nil, fmt.Errorf("Feed did not match expected type of *json.Feed")
	}

	result := &Feed{}
	result.FeedVersion = json.Version
	result.Title = t.convertFeedTitle(json)
	result.Link = t.convertFeedLink(json)
	result.FeedLink = t.convertFeedFeedLink(json)
	result.Links = t.convertFeedLinks(json)
	result.Description = t.convertFeedDescription(json)
	result.Image = t.convertFeedImage(json)
	result.Authors = t.convertFeedAuthors(json)
	result.Language = t.convertFeedLanguage(json)
	result.Items = t.convertFeedItems(json)
	result.Updated = t.convertFeedUpdated(json)
	result.UpdatedParsed = t.convertFeedUpdatedParsed(json)
	result.Published = t.convertFeedPublished(json)
	result.PublishedParsed = t.convertFeedPublishedParsed(json)
	result.FeedType = "json"
	// TODO UserComment is missing in global Feed
	// TODO NextURL is missing in global Feed
	// TODO Favicon is missing in global Feed
	// TODO Exipred is missing in global Feed
	// TODO Hubs is not supported in json.Feed
	// TODO Extensions is not supported in json.Feed
	return result, nil
}

func (t *DefaultJSONConverter) convertFeedItem(jsonItem *json.Item) (item *Item) {
	item = &Item{}
	item.GUID = t.convertItemGUID(jsonItem)
	item.Link = t.convertItemLink(jsonItem)
	item.Links = t.convertItemLinks(jsonItem)
	item.Title = t.convertItemTitle(jsonItem)
	item.Content = t.convertItemContent(jsonItem)
	item.Description = t.convertItemDescription(jsonItem)
	item.Image = t.convertItemImage(jsonItem)
	item.Published = t.convertItemPublished(jsonItem)
	item.PublishedParsed = t.convertItemPublishedParsed(jsonItem)
	item.Updated = t.convertItemUpdated(jsonItem)
	item.UpdatedParsed = t.convertItemUpdatedParsed(jsonItem)
	item.Authors = t.convertItemAuthors(jsonItem)
	item.Categories = t.convertItemCategories(jsonItem)
	item.Enclosures = t.convertItemEnclosures(jsonItem)
	// TODO ExternalURL is missing in global Feed
	// TODO BannerImage is missing in global Feed
	return
}

func (t *DefaultJSONConverter) convertFeedTitle(json *json.Feed) (title string) {
	if json.Title != "" {
		title = json.Title
	}
	return
}

func (t *DefaultJSONConverter) convertFeedDescription(json *json.Feed) (desc string) {
	return json.Description
}

func (t *DefaultJSONConverter) convertFeedLink(json *json.Feed) (link string) {
	if json.HomePageURL != "" {
		link = json.HomePageURL
	}
	return
}

func (t *DefaultJSONConverter) convertFeedFeedLink(json *json.Feed) (link string) {
	if json.FeedURL != "" {
		link = json.FeedURL
	}
	return
}

func (t *DefaultJSONConverter) convertFeedLinks(json *json.Feed) (links []string) {
	if json.HomePageURL != "" {
		links = append(links, json.HomePageURL)
	}
	if json.FeedURL != "" {
		links = append(links, json.FeedURL)
	}
	return
}

func (t *DefaultJSONConverter) convertFeedUpdated(json *json.Feed) (updated string) {
	if len(json.Items) > 0 {
		updated = json.Items[0].DateModified
	}
	return
}

func (t *DefaultJSONConverter) convertFeedUpdatedParsed(json *json.Feed) (updated *time.Time) {
	if len(json.Items) > 0 {
		updateTime, err := shared.ParseDate(json.Items[0].DateModified)
		if err == nil {
			updated = &updateTime
		}
	}
	return
}

func (t *DefaultJSONConverter) convertFeedPublished(json *json.Feed) (published string) {
	if len(json.Items) > 0 {
		published = json.Items[0].DatePublished
	}
	return
}

func (t *DefaultJSONConverter) convertFeedPublishedParsed(json *json.Feed) (published *time.Time) {
	if len(json.Items) > 0 {
		publishTime, err := shared.ParseDate(json.Items[0].DatePublished)
		if err == nil {
			published = &publishTime
		}
	}
	return
}

func (t *DefaultJSONConverter) convertFeedAuthor(json *json.Feed) (author *Person) {
	if json.Author != nil {
		name, address := shared.ParseNameAddress(json.Author.Name)
		author = &Person{}
		author.Name = name
		author.Email = address
	}
	// Author.URL is missing in global feed
	// Author.Avatar is missing in global feed
	return
}

func (t *DefaultJSONConverter) convertFeedAuthors(json *json.Feed) (authors []*Person) {
	if json.Authors != nil {
		authors = []*Person{}
		for _, a := range json.Authors {
			name, address := shared.ParseNameAddress(a.Name)
			author := &Person{}
			author.Name = name
			author.Email = address

			authors = append(authors, author)
		}
	} else if author := t.convertFeedAuthor(json); author != nil {
		authors = []*Person{author}
	}
	// Author.URL is missing in global feed
	// Author.Avatar is missing in global feed
	return
}

func (t *DefaultJSONConverter) convertFeedLanguage(json *json.Feed) (language string) {
	language = json.Language
	return
}

func (t *DefaultJSONConverter) convertFeedImage(json *json.Feed) (image *Image) {
	// Using the Icon rather than the image
	// icon (optional, string) is the URL of an image for the feed suitable to be used in a timeline. It should be square and relatively large — such as 512 x 512
	if json.Icon != "" {
		image = &Image{}
		image.URL = json.Icon
	}
	return
}

func (t *DefaultJSONConverter) convertFeedItems(json *json.Feed) (items []*Item) {
	items = []*Item{}
	for _, i := range json.Items {
		items = append(items, t.convertFeedItem(i))
	}
	return
}

func (t *DefaultJSONConverter) convertItemTitle(jsonItem *json.Item) (title string) {
	if jsonItem.Title != "" {
		title = jsonItem.Title
	}
	return
}

func (t *DefaultJSONConverter) convertItemDescription(jsonItem *json.Item) (desc string) {
	if jsonItem.Summary != "" {
		desc = jsonItem.Summary
	}
	return
}

func (t *DefaultJSONConverter) convertItemContent(jsonItem *json.Item) (content string) {
	if jsonItem.ContentHTML != "" {
		content = jsonItem.ContentHTML
	} else if jsonItem.ContentText != "" {
		content = jsonItem.ContentText
	}
	return
}

func (t *DefaultJSONConverter) convertItemLink(jsonItem *json.Item) (link string) {
	return jsonItem.URL
}

func (t *DefaultJSONConverter) convertItemLinks(jsonItem *json.Item) (links []string) {
	if jsonItem.URL != "" {
		links = append(links, jsonItem.URL)
	}
	if jsonItem.ExternalURL != "" {
		links = append(links, jsonItem.ExternalURL)
	}
	return
}

func (t *DefaultJSONConverter) convertItemUpdated(jsonItem *json.Item) (updated string) {
	if jsonItem.DateModified != "" {
		updated = jsonItem.DateModified
	}
	return updated
}

func (t *DefaultJSONConverter) convertItemUpdatedParsed(jsonItem *json.Item) (updated *time.Time) {
	if jsonItem.DateModified != "" {
		updatedTime, err := shared.ParseDate(jsonItem.DateModified)
		if err == nil {
			updated = &updatedTime
		}
	}
	return
}

func (t *DefaultJSONConverter) convertItemPublished(jsonItem *json.Item) (pubDate string) {
	if jsonItem.DatePublished != "" {
		pubDate = jsonItem.DatePublished
	}
	return
}

func (t *DefaultJSONConverter) convertItemPublishedParsed(jsonItem *json.Item) (pubDate *time.Time) {
	if jsonItem.DatePublished != "" {
		publishTime, err := shared.ParseDate(jsonItem.DatePublished)
		if err == nil {
			pubDate = &publishTime
		}
	}
	return
}

func (t *DefaultJSONConverter) convertItemAuthor(jsonItem *json.Item) (author *Person) {
	if jsonItem.Author != nil {
		name, address := shared.ParseNameAddress(jsonItem.Author.Name)
		author = &Person{}
		author.Name = name
		author.Email = address
	}
	// Author.URL is missing in global feed
	// Author.Avatar is missing in global feed
	return
}

func (t *DefaultJSONConverter) convertItemAuthors(jsonItem *json.Item) (authors []*Person) {
	if jsonItem.Authors != nil {
		authors = []*Person{}
		for _, a := range jsonItem.Authors {
			name, address := shared.ParseNameAddress(a.Name)
			author := &Person{}
			author.Name = name
			author.Email = address

			authors = append(authors, author)
		}
	} else if author := t.convertItemAuthor(jsonItem); author != nil {
		authors = []*Person{author}
	}
	// Author.URL is missing in global feed
	// Author.Avatar is missing in global feed
	return
}

func (t *DefaultJSONConverter) convertItemGUID(jsonItem *json.Item) (guid string) {
	if jsonItem.ID != "" {
		guid = jsonItem.ID
	}
	return
}

func (t *DefaultJSONConverter) convertItemImage(jsonItem *json.Item) (image *Image) {
	if jsonItem.Image != "" {
		image = &Image{}
		image.URL = jsonItem.Image
	} else if jsonItem.BannerImage != "" {
		image = &Image{}
		image.URL = jsonItem.BannerImage
	}
	return
}

func (t *DefaultJSONConverter) convertItemCategories(jsonItem *json.Item) (categories []string) {
	if len(jsonItem.Tags) > 0 {
		categories = jsonItem.Tags
	}
	return
}

func (t *DefaultJSONConverter) convertItemEnclosures(jsonItem *json.Item) (enclosures []*Enclosure) {
	if jsonItem.Attachments != nil {
		for _, attachment := range *jsonItem.Attachments {
			e := &Enclosure{}
			e.URL = attachment.URL
			e.Type = attachment.MimeType
			e.Length = fmt.Sprintf("%d", attachment.DurationInSeconds)
			// Title is not defined in global enclosure
			// SizeInBytes is not defined in global enclosure
			enclosures = append(enclosures, e)
		}
	}
	return
}
