package assert

type v2Pattern struct {
	placeholders    map[string]v2Placeholder
	items           []v2Item
	trailingNewline bool
}

type v2Placeholder struct {
	Name     string
	Type     string
	Metadata map[string]string
}

type v2Item interface {
	isV2Item()
}

type v2LiteralLine struct {
	Text string
}

func (v2LiteralLine) isV2Item() {}

type v2PatternLine struct {
	Segments []v2Segment
}

func (v2PatternLine) isV2Item() {}

type v2RegexLine struct {
	Pattern string
}

func (v2RegexLine) isV2Item() {}

type v2OmitLine struct {
	Count int
}

func (v2OmitLine) isV2Item() {}

type v2Segment interface {
	isV2Segment()
}

type v2Literal struct {
	Text string
}

func (v2Literal) isV2Segment() {}

type v2PlaceholderRef struct {
	Name string
}

func (v2PlaceholderRef) isV2Segment() {}

type v2Color struct {
	Spec colorSpec
	Text string
}

func (v2Color) isV2Segment() {}

type colorSpec struct {
	Tokens []string
}