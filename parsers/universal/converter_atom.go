package universal

import (
	"fmt"
	"strings"
	"time"

	"github.com/mmcdole/gofeed/v2/parsers/atom"
)

// DefaultAtomConverter converts an atom.Feed struct
// into the generic Feed struct.
//
// This default implementation defines a set of
// mapping rules between atom.Feed -> Feed
// for each of the fields in Feed.
type DefaultAtomConverter struct{}

// convert converts an Atom feed into the universal
// feed type.
func (t *DefaultAtomConverter) Convert(feed interface{}) (*Feed, error) {
	atom, found := feed.(*atom.Feed)
	if !found {
		return nil, fmt.Errorf("Feed did not match expected type of *atom.Feed")
	}

	result := &Feed{}
	result.Title = t.convertFeedTitle(atom)
	result.Description = t.convertFeedDescription(atom)
	result.Link = t.convertFeedLink(atom)
	result.FeedLink = t.convertFeedFeedLink(atom)
	result.Links = t.convertFeedLinks(atom)
	result.Updated = t.convertFeedUpdated(atom)
	result.UpdatedParsed = t.convertFeedUpdatedParsed(atom)
	result.Authors = t.convertFeedAuthors(atom)
	result.Language = t.convertFeedLanguage(atom)
	result.Image = t.convertFeedImage(atom)
	result.Copyright = t.convertFeedCopyright(atom)
	result.Categories = t.convertFeedCategories(atom)
	result.Generator = t.convertFeedGenerator(atom)
	result.Items = t.convertFeedItems(atom)
	result.Extensions = atom.Extensions
	result.FeedVersion = atom.Version
	result.FeedType = "atom"
	return result, nil
}

func (t *DefaultAtomConverter) convertFeedItem(entry *atom.Entry) (item *Item) {
	item = &Item{}
	item.Title = t.convertItemTitle(entry)
	item.Description = t.convertItemDescription(entry)
	item.Content = t.convertItemContent(entry)
	item.Link = t.convertItemLink(entry)
	item.Links = t.convertItemLinks(entry)
	item.Updated = t.convertItemUpdated(entry)
	item.UpdatedParsed = t.convertItemUpdatedParsed(entry)
	item.Published = t.convertItemPublished(entry)
	item.PublishedParsed = t.convertItemPublishedParsed(entry)
	item.Authors = t.convertItemAuthors(entry)
	item.GUID = t.convertItemGUID(entry)
	item.Image = t.convertItemImage(entry)
	item.Categories = t.convertItemCategories(entry)
	item.Enclosures = t.convertItemEnclosures(entry)
	item.Extensions = entry.Extensions
	return
}

func (t *DefaultAtomConverter) convertFeedTitle(atom *atom.Feed) (title string) {
	return atom.Title
}

func (t *DefaultAtomConverter) convertFeedDescription(atom *atom.Feed) (desc string) {
	return atom.Subtitle
}

func (t *DefaultAtomConverter) convertFeedLink(atom *atom.Feed) (link string) {
	l := t.firstLinkWithType("alternate", atom.Links)
	if l != nil {
		link = l.Href
	}
	return
}

func (t *DefaultAtomConverter) convertFeedFeedLink(atom *atom.Feed) (link string) {
	feedLink := t.firstLinkWithType("self", atom.Links)
	if feedLink != nil {
		link = feedLink.Href
	}
	return
}

func (t *DefaultAtomConverter) convertFeedLinks(atom *atom.Feed) (links []string) {
	for _, l := range atom.Links {
		if l.Rel == "" || l.Rel == "alternate" || l.Rel == "self" {
			links = append(links, l.Href)
		}
	}
	return
}

func (t *DefaultAtomConverter) convertFeedUpdated(atom *atom.Feed) (updated string) {
	return atom.Updated
}

func (t *DefaultAtomConverter) convertFeedUpdatedParsed(atom *atom.Feed) (updated *time.Time) {
	return atom.UpdatedParsed
}

func (t *DefaultAtomConverter) convertFeedAuthor(atom *atom.Feed) (author *Person) {
	a := t.firstPerson(atom.Authors)
	if a != nil {
		feedAuthor := Person{}
		feedAuthor.Name = a.Name
		feedAuthor.Email = a.Email
		author = &feedAuthor
	}
	return
}

func (t *DefaultAtomConverter) convertFeedAuthors(atom *atom.Feed) (authors []*Person) {
	if atom.Authors != nil {
		authors = []*Person{}

		for _, a := range atom.Authors {
			authors = append(authors, &Person{
				Name:  a.Name,
				Email: a.Email,
			})
		}
	}

	return
}

func (t *DefaultAtomConverter) convertFeedLanguage(atom *atom.Feed) (language string) {
	return atom.Language
}

func (t *DefaultAtomConverter) convertFeedImage(atom *atom.Feed) (image *Image) {
	if atom.Logo != "" {
		feedImage := Image{}
		feedImage.URL = atom.Logo
		image = &feedImage
	}
	return
}

func (t *DefaultAtomConverter) convertFeedCopyright(atom *atom.Feed) (rights string) {
	return atom.Rights
}

func (t *DefaultAtomConverter) convertFeedGenerator(atom *atom.Feed) (generator string) {
	if atom.Generator != nil {
		if atom.Generator.Value != "" {
			generator += atom.Generator.Value
		}
		if atom.Generator.Version != "" {
			generator += " v" + atom.Generator.Version
		}
		if atom.Generator.URI != "" {
			generator += " " + atom.Generator.URI
		}
		generator = strings.TrimSpace(generator)
	}
	return
}

func (t *DefaultAtomConverter) convertFeedCategories(atom *atom.Feed) (categories []string) {
	if atom.Categories != nil {
		categories = []string{}
		for _, c := range atom.Categories {
			categories = append(categories, c.Term)
		}
	}
	return
}

func (t *DefaultAtomConverter) convertFeedItems(atom *atom.Feed) (items []*Item) {
	items = []*Item{}
	for _, entry := range atom.Entries {
		items = append(items, t.convertFeedItem(entry))
	}
	return
}

func (t *DefaultAtomConverter) convertItemTitle(entry *atom.Entry) (title string) {
	return entry.Title
}

func (t *DefaultAtomConverter) convertItemDescription(entry *atom.Entry) (desc string) {
	return entry.Summary
}

func (t *DefaultAtomConverter) convertItemContent(entry *atom.Entry) (content string) {
	if entry.Content != nil {
		content = entry.Content.Value
	}
	return
}

func (t *DefaultAtomConverter) convertItemLink(entry *atom.Entry) (link string) {
	l := t.firstLinkWithType("alternate", entry.Links)
	if l != nil {
		link = l.Href
	}
	return
}

func (t *DefaultAtomConverter) convertItemLinks(entry *atom.Entry) (links []string) {
	for _, l := range entry.Links {
		if l.Rel == "" || l.Rel == "alternate" || l.Rel == "self" {
			links = append(links, l.Href)
		}
	}
	return
}

func (t *DefaultAtomConverter) convertItemUpdated(entry *atom.Entry) (updated string) {
	return entry.Updated
}

func (t *DefaultAtomConverter) convertItemUpdatedParsed(entry *atom.Entry) (updated *time.Time) {
	return entry.UpdatedParsed
}

func (t *DefaultAtomConverter) convertItemPublished(entry *atom.Entry) (published string) {
	published = entry.Published
	if published == "" {
		published = entry.Updated
	}
	return
}

func (t *DefaultAtomConverter) convertItemPublishedParsed(entry *atom.Entry) (published *time.Time) {
	published = entry.PublishedParsed
	if published == nil {
		published = entry.UpdatedParsed
	}
	return
}

func (t *DefaultAtomConverter) convertItemAuthor(entry *atom.Entry) (author *Person) {
	a := t.firstPerson(entry.Authors)
	if a != nil {
		author = &Person{}
		author.Name = a.Name
		author.Email = a.Email
	}
	return
}

func (t *DefaultAtomConverter) convertItemAuthors(entry *atom.Entry) (authors []*Person) {
	if entry.Authors != nil {
		authors = []*Person{}
		for _, a := range entry.Authors {
			authors = append(authors, &Person{
				Name:  a.Name,
				Email: a.Email,
			})
		}
	}
	return
}

func (t *DefaultAtomConverter) convertItemGUID(entry *atom.Entry) (guid string) {
	return entry.ID
}

func (t *DefaultAtomConverter) convertItemImage(entry *atom.Entry) (image *Image) {
	return nil
}

func (t *DefaultAtomConverter) convertItemCategories(entry *atom.Entry) (categories []string) {
	if entry.Categories != nil {
		categories = []string{}
		for _, c := range entry.Categories {
			categories = append(categories, c.Term)
		}
	}
	return
}

func (t *DefaultAtomConverter) convertItemEnclosures(entry *atom.Entry) (enclosures []*Enclosure) {
	if entry.Links != nil {
		enclosures = []*Enclosure{}
		for _, e := range entry.Links {
			if e.Rel == "enclosure" {
				enclosure := &Enclosure{}
				enclosure.URL = e.Href
				enclosure.Length = e.Length
				enclosure.Type = e.Type
				enclosures = append(enclosures, enclosure)
			}
		}

		if len(enclosures) == 0 {
			enclosures = nil
		}
	}
	return
}

func (t *DefaultAtomConverter) firstLinkWithType(linkType string, links []*atom.Link) *atom.Link {
	if links == nil {
		return nil
	}

	for _, link := range links {
		if link.Rel == linkType {
			return link
		}
	}
	return nil
}

func (t *DefaultAtomConverter) firstPerson(persons []*atom.Person) (person *atom.Person) {
	if persons == nil || len(persons) == 0 {
		return
	}

	person = persons[0]
	return
}
