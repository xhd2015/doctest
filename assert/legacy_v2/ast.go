package legacy_v2

// Pattern is an immutable parsed v2 output template.
type Pattern struct {
	placeholders    map[string]placeholder
	items           []item
	trailingNewline bool
}

type placeholder struct {
	Name     string
	Type     string
	Metadata map[string]string
}

type item interface {
	isItem()
}

type literalLine struct {
	Text string
}

func (literalLine) isItem() {}

type patternLine struct {
	Segments []segment
}

func (patternLine) isItem() {}

type regexLine struct {
	Pattern string
}

func (regexLine) isItem() {}

type omitLine struct {
	Count int
}

func (omitLine) isItem() {}

type segment interface {
	isSegment()
}

type literal struct {
	Text string
}

func (literal) isSegment() {}

type placeholderRef struct {
	Name string
}

func (placeholderRef) isSegment() {}

type color struct {
	Spec colorSpec
	Text string
}

func (color) isSegment() {}

type colorSpec struct {
	Tokens []string
}
