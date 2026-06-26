# Scenario

**Feature**: Output Assert DSL parse and match harness

```
# author supplies template + actual CLI output
Author -> Parser: template text
Parser -> Matcher: immutable Pattern AST
Matcher <- actual output bytes (+ options)
Matcher -> pass or rich diff error
```

## Preconditions

- Package `github.com/xhd2015/doctest/assert` exposes `Parse`, `Match`, `Contains`, and `NormalizeNewlines`.
- Root helpers summarize parsed AST shape for parse-only leaves.

## Steps

1. Ancestor `Setup` functions set `req.Operation`, `req.Template`, `req.Actual`, and `req.Options`.
2. Root `Run` parses the template; parse-only leaves stop after summarizing the AST.
3. Match leaves call `assert.Match` with options and return match errors verbatim.

## Context

- `ParsedSummary` uses a stable `kind+details` DSL (e.g. `LiteralLine×2`, `BlockOptional{1}`) that parse leaves assert against.
- Match leaves assert `resp.MatchOK` or expected error substrings in `resp.MatchError`.

```go
import (
	"fmt"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/assert"
)

func summarizePattern(p assert.Pattern) string {
	return assert.PatternSummary(p)
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

func greenWrap(s string) string {
	return fmt.Sprintf("\x1b[32m%s\x1b[0m", s)
}

func boldWrap(s string) string {
	return fmt.Sprintf("\x1b[1m%s\x1b[0m", s)
}

func boldGrayWrap(s string) string {
	return fmt.Sprintf("\x1b[1m\x1b[90m%s\x1b[0m", s)
}

func raw256Wrap(s string) string {
	return fmt.Sprintf("\x1b[38;5;208m%s\x1b[0m", s)
}
```