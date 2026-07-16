# Scenario

**Feature**: Output Assert DSL v3 full suite (migrated cookbook + match/parse leaves)

```
# author supplies v3 YAML header + strict template body (raw RE content lines)
Author -> Facade.Parse: template text
Facade -> v3 Parser (version: 3 or dialect-no-version) / legacy_v2 / legacy_v1
Parser -> Matcher: immutable Pattern AST
Matcher <- actual output bytes
Matcher -> pass or rich error (bindings, line mismatch)
```

## Preconditions

- Package `github.com/xhd2015/doctest/assert` exposes `Parse`, `Match`, and `PatternSummary`.
- Facade routes:
  - `version: 3` or YAML dialect with placeholders and no version → v3
  - `version: 2` → legacy_v2 (only thin smoke elsewhere; this suite is v3-primary)
  - pure v1 tags / non-dialect frontmatter → legacy_v1
- v3 content lines are **raw Go regex** full-line matches (`^…$`). Authors escape
  RE metacharacters for literals (`.` → `\.`, `(` → `\(`, …). Color-tag **inner**
  text is QuoteMeta'd by the engine — no escape needed inside `<ansi-color>…</ansi-color>`.
- Placeholders: `type=string` / `type=number` / `regex=`; repeated `__NAME__` same-value binding.
- No `MatchContains` mode — strict full match only.

## Steps

1. Ancestor `Setup` functions set `req.Operation`, `req.Template`, `req.Actual`, and flags.
2. Root `Run` calls `assert.Parse` then `assert.Match` for match leaves.
3. Parse leaves summarize AST via `assert.PatternSummary`.

## Context

- `ParsedSummary` uses stable v3 shapes (e.g. `Placeholders:PORT{number}`, `RegexLine`,
  `OmitLine{3}`). Content lines always summarize as `RegexLine` (no separate LiteralLine).
- v1 fallback leaves assert `ContainsBlock`, `LiteralLine`, or other legacy_v1 shapes.
- Color helpers produce ANSI bytes for `<ansi-color>` span matching (same SGR tokens as v2).
- **P3 mapping:** this tree was migrated from `output-assert-v2` (legacy_v2 templates)
  to v3 defaults — see root `DOCTEST.md` Mapping notes.

```go
import (
	"fmt"
	"strings"
	"testing"
)

// v3Template builds an explicit version: 3 document (canonical for this suite).
func v3Template(header, body string) string {
	return "---\nversion: 3\n" + header + "---\n" + body
}

// In ASSERT.md Go blocks you may also write the template as one raw string:
//
//	req.Template = `---
//	version: 3
//	---
//	body line`
//
// Leading blank lines before the opening --- are trimmed; trailing blank lines are not.
// Escape RE metas in body literals (e.g. version 1\.0). Color-tag inners stay unescaped.

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

func requireSummary(t *testing.T, resp *Response, want string) {
	t.Helper()
	requireParseOK(t, resp)
	if resp.ParsedSummary != want {
		t.Fatalf("parsed summary:\n  got:  %q\n  want: %q", resp.ParsedSummary, want)
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

func requireSummaryNotContains(t *testing.T, resp *Response, subs ...string) {
	t.Helper()
	requireParseOK(t, resp)
	for _, sub := range subs {
		if strings.Contains(resp.ParsedSummary, sub) {
			t.Fatalf("parsed summary %q unexpectedly contains %q", resp.ParsedSummary, sub)
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

func greenBoldWrap(s string) string {
	return fmt.Sprintf("\x1b[1m\x1b[32m%s\x1b[0m", s)
}
```
