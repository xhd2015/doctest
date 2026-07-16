# Scenario

**Feature**: Output Assert DSL v3 parse and match harness

```
# author supplies v3 YAML header + strict template body (raw RE content lines)
Author -> Facade.Parse: template text
Facade -> v3 Parser (version: 3 or dialect-no-version) or legacy_v2 / legacy_v1
Parser -> Matcher: immutable Pattern AST
Matcher <- actual output bytes
Matcher -> pass or rich error (bindings, line mismatch)
```

## Preconditions

- Package `github.com/xhd2015/doctest/assert` exposes `Parse`, `Match`, and `PatternSummary`.
- Facade routes:
  - `version: 2` → legacy_v2
  - pure v1 tags (no YAML dialect) → legacy_v1
  - `version: 3` or YAML dialect without version → v3
  - unknown version → parse error
- v3 content lines are raw Go regex full-line matches; omit markers and color tags are special.
- v3 has no `hasRegexIntent` and no `MatchContains` mode — strict full match only.
- Placeholders: `type=string` / `type=number` / `regex=`; dual type+regex is a parse error.
- Repeated `__NAME__` requires same-value binding at match time.

## Steps

1. Ancestor `Setup` functions set `req.Operation`, `req.Template`, `req.Actual`, and flags.
2. Root `Run` calls `assert.Parse` then `assert.Match` for match leaves.
3. Parse leaves summarize AST via `assert.PatternSummary`.

## Context

- `ParsedSummary` should expose stable shapes (e.g. placeholder names/types, `OmitLine{N}`,
  `RegexLine`) once v3 is implemented; leaves use substring checks.
- Color helpers produce ANSI bytes for `<ansi-color>` span matching (same SGR tokens as v2).
- Classic TDD: most v3 leaves fail until the engine lands; routing leaves may already pass.

```go
import (
	"fmt"
	"strings"
	"testing"
)

// v3Template builds an explicit version: 3 document.
func v3Template(header, body string) string {
	return "---\nversion: 3\n" + header + "---\n" + body
}

// v3TemplateNoVersion builds a YAML dialect header without a version key (default → v3).
func v3TemplateNoVersion(header, body string) string {
	return "---\n" + header + "---\n" + body
}

// v2Template builds a version: 2 document (legacy_v2 routing).
func v2Template(header, body string) string {
	return "---\nversion: 2\n" + header + "---\n" + body
}

func requireParseOK(t *testing.T, resp *Response) {
	t.Helper()
	if !resp.ParseOK {
		t.Fatalf("expected parse success, got error: %s", resp.ParseError)
	}
}

func requireParseError(t *testing.T, resp *Response) {
	t.Helper()
	if resp.ParseOK {
		t.Fatalf("expected parse error, got summary: %s", resp.ParsedSummary)
	}
	if resp.ParseError == "" {
		t.Fatal("expected non-empty parse error")
	}
}

func requireSummaryContains(t *testing.T, resp *Response, subs ...string) {
	t.Helper()
	requireParseOK(t, resp)
	for _, sub := range subs {
		if !strings.Contains(resp.ParsedSummary, sub) {
			t.Fatalf("parsed summary %q missing %q", resp.ParsedSummary, sub)
		}
	}
}

func requireMatchOK(t *testing.T, resp *Response) {
	t.Helper()
	if !resp.MatchOK {
		t.Fatalf("expected match success, got: %s", resp.MatchError)
	}
}

func requireMatchError(t *testing.T, resp *Response, subs ...string) {
	t.Helper()
	if resp.MatchOK {
		t.Fatal("expected match failure, got success")
	}
	for _, sub := range subs {
		if !strings.Contains(resp.MatchError, sub) {
			t.Fatalf("match error %q missing %q", resp.MatchError, sub)
		}
	}
}

func requireParseErrorContains(t *testing.T, resp *Response, subs ...string) {
	t.Helper()
	requireParseError(t, resp)
	for _, sub := range subs {
		if !strings.Contains(resp.ParseError, sub) {
			t.Fatalf("parse error %q missing %q", resp.ParseError, sub)
		}
	}
}

func grayWrap(s string) string {
	return fmt.Sprintf("\x1b[90m%s\x1b[0m", s)
}
```
