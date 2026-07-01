package assert

type Pattern struct {
	items           []Item
	trailingNewline bool
}

type Item interface {
	isItem()
}

type LiteralLine struct {
	Text string
}

func (LiteralLine) isItem() {}

type PatternLine struct {
	Segments []Segment
}

func (PatternLine) isItem() {}

type Segment interface {
	isSegment()
}

type Literal struct {
	Text string
}

func (Literal) isSegment() {}

type Hint struct {
	Label string
	Text  string
}

func (Hint) isSegment() {}

type AnsiColor struct {
	Spec AnsiColorSpec
	Text string
}

func (AnsiColor) isSegment() {}

type AnsiColorSpec struct {
	Tokens []string
}

type InlineOptional struct {
	Text string
}

func (InlineOptional) isSegment() {}

type InlineAnyOf struct {
	Branches []InlineExpectBranch
}

func (InlineAnyOf) isSegment() {}

type InlineExpectBranch struct {
	Segments []Segment
}

type InlineRegex struct {
	Pattern string
}

func (InlineRegex) isSegment() {}

type RegexLine struct {
	Pattern string
}

func (RegexLine) isItem() {}

type BlockOptional struct {
	Items []Item
}

func (BlockOptional) isItem() {}

type AnyOfBlock struct {
	Branches []ExpectBranch
}

func (AnyOfBlock) isItem() {}

type ExpectBranch struct {
	Items []Item
}

type ContainsBlock struct {
	Fragments []ContainsFragment
}

func (ContainsBlock) isItem() {}

type ContainsMatchMode int

const (
	ContainsSubstring ContainsMatchMode = iota
	ContainsStartWith
	ContainsEndWith
)

type ContainsFragment struct {
	Mode     ContainsMatchMode
	Text     string
	Segments []Segment
}
