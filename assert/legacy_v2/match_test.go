package legacy_v2

import "testing"

func headerBody(header, body string) (string, string) {
	return "version: 2\n" + header, body
}

func TestMatchStringPlaceholderEOL(t *testing.T) {
	// Whole-line string placeholder (real-world cookbook shape).
	h, b := headerBody("__LINE__: type=string\n", "__LINE__")
	p, err := Parse(h, b)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if err := Match(p, "[10/10] Linking CXX executable app", true); err != nil {
		t.Fatalf("match: %v", err)
	}
}

func TestMatchStringPlaceholderMidLine(t *testing.T) {
	// V2-M3 style: literal + string placeholder.
	h, b := headerBody("__USER__: type=string\n", "Hello __USER__")
	p, err := Parse(h, b)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if err := Match(p, "Hello alice", true); err != nil {
		t.Fatalf("match: %v", err)
	}
}

func TestMatchStringPlaceholderBetweenLiterals(t *testing.T) {
	h, b := headerBody("__X__: type=string\n", "pre __X__ post")
	p, err := Parse(h, b)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if err := Match(p, "pre middle post", true); err != nil {
		t.Fatalf("match: %v", err)
	}
	if err := Match(p, "pre wrong", true); err == nil {
		t.Fatal("expected mismatch when trailing literal missing")
	}
}

func TestMatchNumberStillWorks(t *testing.T) {
	h, b := headerBody("__PORT__: type=number\n", "Server listen on: __PORT__")
	p, err := Parse(h, b)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if err := Match(p, "Server listen on: 8901", true); err != nil {
		t.Fatalf("match: %v", err)
	}
}

func TestMatchRegexLinePass(t *testing.T) {
	// V2-M4: pure regex-intent line must not QuoteMeta the body.
	h, b := headerBody("", ".*Some middle content.*suffix content")
	p, err := Parse(h, b)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if err := Match(p, "XXXSome middle contentYYYsuffix content", true); err != nil {
		t.Fatalf("match: %v", err)
	}
}

func TestMatchRegexAlternation(t *testing.T) {
	// V2-M13
	h, b := headerBody("", "(ok|fail)")
	p, err := Parse(h, b)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if err := Match(p, "ok", true); err != nil {
		t.Fatalf("match: %v", err)
	}
}
