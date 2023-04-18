package config

type ParseOptions struct {
	// MaxItems specifies the maximum number of feed items to parse.
	// The default value is 0, which means no limit.
	MaxItems int

	// ParseDates determines if the feed parser will attempt to parse dates into `time.Time` objects.
	// The default value is true.
	ParseDates bool

	// ParseExtensions determines if the feed parser will attempt to parse feed extensions such as
	// iTunes extensions, or custom feed extensions.
	// The default value is true.
	ParseExtensions bool

	// KeepOriginalFeed specifies if the parser should retain the raw feed in the `Feed` struct's `RawFeed` field.
	// The default value is false.
	KeepOriginalFeed bool

	// StrictnessOptions holds the options for controlling the strictness of the parsing.
	// The defaults are set to their least strict values.
	StrictnessOptions StrictnessOptions
}

type StrictnessOptions struct {
	// StripInvalidCharacters specifies if invalid feed characters should be stripped out.
	// The default value is true.
	StripInvalidCharacters bool

	// AutoCloseTags specifies if the parser should automatically close unclosed tags.
	// The default value is true.
	AutoCloseTags bool

	// AllowUndisclosedXMLNamespaces specifies if the parser should allow undisclosed XML namespaces.
	// The default value is true.
	AllowUndisclosedXMLNamespaces bool

	// AllowCustomXMLElements specifies if the parser should allow custom XML elements.
	// The default value is true.
	AllowCustomXMLElements bool

	// AllowIncorrectDateFormats specifies if the parser should allow incorrect date formats.
	// The default value is true.
	AllowIncorrectDateFormats bool

	// AllowUnescapedMarkup specifies if the parser should allow unescaped / naked markup in feed elements.
	// The default value is true.
	AllowUnescapedMarkup bool
}

func NewParseOptions() ParseOptions {
	return ParseOptions{
		MaxItems:         0,
		ParseDates:       true,
		ParseExtensions:  true,
		KeepOriginalFeed: false,
		StrictnessOptions: StrictnessOptions{
			StripInvalidCharacters:        true,
			AutoCloseTags:                 true,
			AllowUndisclosedXMLNamespaces: true,
			AllowCustomXMLElements:        true,
			AllowIncorrectDateFormats:     true,
			AllowUnescapedMarkup:          true,
		},
	}
}
