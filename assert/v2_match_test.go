package assert

import "testing"

func v2Template(header, body string) string {
	return "---\nversion: 2\n" + header + "---\n" + body
}

func TestV2MatchStringPlaceholderEOL(t *testing.T) {
	// Whole-line string placeholder (real-world cookbook shape).
	tpl := v2Template("__LINE__: type=string\n", "__LINE__")
	p, err := Parse(tpl)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if err := Match(p, "[10/10] Linking CXX executable app"); err != nil {
		t.Fatalf("match: %v", err)
	}
}

func TestV2MatchStringPlaceholderMidLine(t *testing.T) {
	// V2-M3 style: literal + string placeholder.
	tpl := v2Template("__USER__: type=string\n", "Hello __USER__")
	p, err := Parse(tpl)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if err := Match(p, "Hello alice"); err != nil {
		t.Fatalf("match: %v", err)
	}
}

func TestV2MatchStringPlaceholderBetweenLiterals(t *testing.T) {
	tpl := v2Template("__X__: type=string\n", "pre __X__ post")
	p, err := Parse(tpl)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if err := Match(p, "pre middle post"); err != nil {
		t.Fatalf("match: %v", err)
	}
	if err := Match(p, "pre wrong"); err == nil {
		t.Fatal("expected mismatch when trailing literal missing")
	}
}

func TestV2MatchNumberStillWorks(t *testing.T) {
	tpl := v2Template("__PORT__: type=number\n", "Server listen on: __PORT__")
	p, err := Parse(tpl)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if err := Match(p, "Server listen on: 8901"); err != nil {
		t.Fatalf("match: %v", err)
	}
}

func TestV2MatchRegexLinePass(t *testing.T) {
	// V2-M4: pure regex-intent line must not QuoteMeta the body.
	tpl := v2Template("", ".*Some middle content.*suffix content")
	p, err := Parse(tpl)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if err := Match(p, "XXXSome middle contentYYYsuffix content"); err != nil {
		t.Fatalf("match: %v", err)
	}
}

func TestV2MatchRegexAlternation(t *testing.T) {
	// V2-M13
	tpl := v2Template("", "(ok|fail)")
	p, err := Parse(tpl)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if err := Match(p, "ok"); err != nil {
		t.Fatalf("match: %v", err)
	}
}
