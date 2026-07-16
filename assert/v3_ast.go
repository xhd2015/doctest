package assert

// v3Pattern is an immutable parsed v3 output template.
// Content lines are raw Go regex full-line matches; omit markers are special.
type v3Pattern struct {
	placeholders    map[string]v3Placeholder
	items           []v3Item
	trailingNewline bool
}

type v3Placeholder struct {
	Name     string
	Type     string // "string", "number", or "regex"
	Regex    string // custom fragment when Type == "regex"
	Metadata map[string]string
}

type v3Item interface {
	isV3Item()
}

type v3RegexLine struct {
	// Pattern is a full-line RE including ^ and $ anchors, with named groups for placeholders.
	Pattern string
}

func (v3RegexLine) isV3Item() {}

type v3OmitLine struct {
	Count int
}

func (v3OmitLine) isV3Item() {}
