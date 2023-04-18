package extensions

import (
	"strings"

	xpp "github.com/mmcdole/goxpp"
)

type Extension struct {
	Name     string                 `json:"name"`
	Value    string                 `json:"value"`
	Attrs    map[string]string      `json:"attrs"`
	Children map[string][]Extension `json:"children"`
}

// Extensions is the generic extension map for Feeds and Items.
// The first map is for the element namespace prefix (e.g., itunes).
// The second map is for the element name (e.g., author).
type Extensions map[string]map[string][]Extension

// ParseExtension parses the current element of the XMLPullParser as an extension element
// and updates the Extensions map with the parsed Extension.
func ParseExtension(fe Extensions, p *xpp.XMLPullParser) (Extensions, error) {
	// Get the namespace prefix for the current element
	prefix := prefixForNamespace(p.Space, p)
	if prefix == "" {
		prefix = "default"
	}

	// Parse the extension element
	e, err := parseExtensionElement(p, prefix)
	if err != nil {
		return nil, err
	}

	// Ensure the extension prefix map exists
	if _, ok := fe[prefix]; !ok {
		fe[prefix] = make(map[string][]Extension)
	}

	// Ensure the extension element slice exists
	if _, ok := fe[prefix][p.Name]; !ok {
		fe[prefix][p.Name] = []Extension{}
	}

	// Add the parsed extension to the Extensions map
	fe[prefix][p.Name] = append(fe[prefix][p.Name], e)
	return fe, nil
}

// parseExtensionElement parses the current element of the XMLPullParser as an extension element
// and returns the parsed Extension.
func parseExtensionElement(p *xpp.XMLPullParser, prefix string) (e Extension, err error) {
	// Expect the start of an XML tag
	if err = p.Expect(xpp.StartTag, "*"); err != nil {
		return e, err
	}

	// Initialize a new Extension with the current element's name and empty maps for attributes and children
	e.Name = p.Name
	e.Children = make(map[string][]Extension)
	e.Attrs = make(map[string]string)

	// Add the attributes of the current element to the Attributes map
	for _, attr := range p.Attrs {
		e.Attrs[attr.Name.Local] = attr.Value
	}

	// Loop through the XML tokens to find child elements or text content
	for {
		tok, err := p.Next()
		if err != nil {
			return e, err
		}

		if tok == xpp.EndTag {
			break
		}

		if tok == xpp.StartTag {
			// Parse child extension elements
			child, err := parseExtensionElement(p, prefix)
			if err != nil {
				return e, err
			}

			// Ensure the child extension element slice exists
			if _, ok := e.Children[child.Name]; !ok {
				e.Children[child.Name] = []Extension{}
			}

			// Add the parsed child extension to the Children map
			e.Children[child.Name] = append(e.Children[child.Name], child)
		} else if tok == xpp.Text {
			// Append the text content to the Value field
			e.Value += p.Text
		}
	}

	// Trim any leading/trailing whitespace from the Value field
	e.Value = strings.TrimSpace(e.Value)

	// Expect the end of the current XML tag
	err = p.Expect(xpp.EndTag, e.Name)
	if err != nil {
		return e, err
	}

	return e, nil
}

func prefixForNamespace(space string, p *xpp.XMLPullParser) string {
	if prefix, ok := canonicalNamespaces[space]; ok {
		return prefix
	}

	if prefix, ok := p.Spaces[space]; ok {
		return prefix
	}

	return space
}

// Namespaces taken from github.com/kurtmckee/feedparser
// These are used for determining canonical name space prefixes
// for many of the popular RSS/Atom extensions.
//
// These canonical prefixes override any prefixes used in the feed itself.
var canonicalNamespaces = map[string]string{
	"http://webns.net/mvcb/":                                         "admin",
	"http://purl.org/rss/1.0/modules/aggregation/":                   "ag",
	"http://purl.org/rss/1.0/modules/annotate/":                      "annotate",
	"http://media.tangent.org/rss/1.0/":                              "audio",
	"http://backend.userland.com/blogChannelModule":                  "blogChannel",
	"http://creativecommons.org/ns#license":                          "cc",
	"http://web.resource.org/cc/":                                    "cc",
	"http://cyber.law.harvard.edu/rss/creativeCommonsRssModule.html": "creativeCommons",
	"http://backend.userland.com/creativeCommonsRssModule":           "creativeCommons",
	"http://purl.org/rss/1.0/modules/company":                        "co",
	"http://purl.org/rss/1.0/modules/content/":                       "content",
	"http://my.theinfo.org/changed/1.0/rss/":                         "cp",
	"http://purl.org/dc/elements/1.1/":                               "dc",
	"http://purl.org/dc/terms/":                                      "dcterms",
	"http://purl.org/rss/1.0/modules/email/":                         "email",
	"http://purl.org/rss/1.0/modules/event/":                         "ev",
	"http://rssnamespace.org/feedburner/ext/1.0":                     "feedburner",
	"http://freshmeat.net/rss/fm/":                                   "fm",
	"http://xmlns.com/foaf/0.1/":                                     "foaf",
	"http://www.w3.org/2003/01/geo/wgs84_pos#":                       "geo",
	"http://www.georss.org/georss":                                   "georss",
	"http://www.opengis.net/gml":                                     "gml",
	"http://postneo.com/icbm/":                                       "icbm",
	"http://purl.org/rss/1.0/modules/image/":                         "image",
	"http://www.itunes.com/DTDs/PodCast-1.0.dtd":                     "itunes",
	"http://example.com/DTDs/PodCast-1.0.dtd":                        "itunes",
	"http://purl.org/rss/1.0/modules/link/":                          "l",
	"http://search.yahoo.com/mrss":                                   "media",
	"http://search.yahoo.com/mrss/":                                  "media",
	"http://madskills.com/public/xml/rss/module/pingback/":           "pingback",
	"http://prismstandard.org/namespaces/1.2/basic/":                 "prism",
	"http://www.w3.org/1999/02/22-rdf-syntax-ns#":                    "rdf",
	"http://www.w3.org/2000/01/rdf-schema#":                          "rdfs",
	"http://purl.org/rss/1.0/modules/reference/":                     "ref",
	"http://purl.org/rss/1.0/modules/richequiv/":                     "reqv",
	"http://purl.org/rss/1.0/modules/search/":                        "search",
	"http://purl.org/rss/1.0/modules/slash/":                         "slash",
	"http://schemas.xmlsoap.org/soap/envelope/":                      "soap",
	"http://purl.org/rss/1.0/modules/servicestatus/":                 "ss",
	"http://hacks.benhammersley.com/rss/streaming/":                  "str",
	"http://purl.org/rss/1.0/modules/subscription/":                  "sub",
	"http://purl.org/rss/1.0/modules/syndication/":                   "sy",
	"http://schemas.pocketsoap.com/rss/myDescModule/":                "szf",
	"http://purl.org/rss/1.0/modules/taxonomy/":                      "taxo",
	"http://purl.org/rss/1.0/modules/threading/":                     "thr",
	"http://purl.org/rss/1.0/modules/textinput/":                     "ti",
	"http://madskills.com/public/xml/rss/module/trackback/":          "trackback",
	"http://wellformedweb.org/commentAPI/":                           "wfw",
	"http://purl.org/rss/1.0/modules/wiki/":                          "wiki",
	"http://www.w3.org/1999/xhtml":                                   "xhtml",
	"http://www.w3.org/1999/xlink":                                   "xlink",
	"http://www.w3.org/XML/1998/namespace":                           "xml",
	"http://podlove.org/simple-chapters":                             "psc",
}
