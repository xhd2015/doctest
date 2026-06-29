# Output Assert DSL Tests

## Version
0.0.2

## DSN (Domain Specific Notion)

### Participants

- **Author** — writes a template string mixing literal CLI output with DSL tags
  (`<optional>`, `<any-of>`, `<contains>`, `<regex>`, `<hint:label>`, `<ansi-color>`).
- **Parser** — reads the template and builds an immutable `Pattern` AST. Rejects
  unknown tags, bare `<id>` placeholders, unclosed tags, label mismatches, and
  empty `<ansi-color>`.
- **Matcher** — walks actual output line-by-line (or span-by-span for inline
  constructs), consuming literal text, optional regions, any-of branches, regex
  spans, contains fragments, hints, and ANSI-colored spans. Block meta lines
  (`<optional>`, `<any-of>`, `<contains>`, `<expect>`) never consume actual output.
- **Options** — `MatchExact` (default), `MatchContains()` (contiguous subregion),
  `NormalizeNewlines` (`\r\n` → `\n`, default on). Trailing newline policy is strict.

### Behaviors

- **Parse** — template text → `Pattern` or parse error with position/message.
- **Match** — `Pattern` + actual string + options → nil (pass) or rich diff error
  (line numbers, hint labels, all any-of branches on failure).
- **Output** — testing helper: parse + match + `t.Fatal` on failure.

## Decision Tree

```
output-assert
├── parse/                          Operation = parse
│   ├── valid/                      successful AST shapes (P1–P4, P7–P12, P14, P16–P19)
│   │   ├── literal-lines           P1 — two literal lines
│   │   ├── block-optional          P2 — literal + block optional + literal
│   │   ├── hint-standalone         P3 — pattern line with hint:id
│   │   ├── hint-with-literal       P4 — literal + hint segments
│   │   ├── inline-optional         P7 — inline optional segment
│   │   ├── block-any-of            P8 — AnyOfBlock with two branches
│   │   ├── contains-block          P9 — ContainsBlock with three fragments
│   │   ├── ansi-color-named        P11 — AnsiColor named gray
│   │   ├── contains-start-with     P12 — ContainsFragment StartWith mode
│   │   ├── ansi-color-raw          P14 — AnsiColor raw #38;5;208
│   │   ├── regex-block-inline      P16 — RegexLine + InlineRegex
│   │   ├── inline-any-of           P17 — inline any-of on pattern line
│   │   ├── ansi-color-combined     P18 — bold + gray tokens
│   │   └── escaped-literal         P19 — backslash-escaped tag text
│   └── errors/                     parse must fail (P5, P6, P10, P13, P15)
│       ├── bare-tag                P5 — bare <id> rejected
│       ├── unclosed-hint           P6 — unclosed hint with position
│       ├── hint-label-mismatch     P10 — open/close label mismatch
│       ├── empty-ansi-color        P13 — empty inner text
│       └── unknown-ansi-name       P15 — unknown color name without #
├── match/                          Operation = match (default MatchExact)
│   ├── literal/                    M1–M3
│   │   ├── exact-pass              M1
│   │   ├── strict-trailing-fail    M2 — template without \n, actual with \n
│   │   ├── strict-trailing-pass    M2b — both agree on trailing newline
│   │   └── line-mismatch           M3 — line 2 differs
│   ├── optional/                   O1–O9
│   │   ├── block-absent            O1
│   │   ├── block-present           O2
│   │   ├── block-missing-newline   O3
│   │   ├── inline-absent           O4
│   │   ├── inline-present          O5
│   │   ├── inline-prefix-required  O6
│   │   ├── meta-lines-ignored      O7
│   │   ├── adjacent-blocks         O8 — separate semantics, line2 only
│   │   └── partial-inner-rejected  O9 — cannot match subset of inner lines
│   ├── any-of/                     A1–A5
│   │   ├── first-branch            A1
│   │   ├── second-branch           A2
│   │   ├── no-branch               A3 — reports all branches
│   │   ├── multiline-branch        A4
│   │   └── meta-lines-ignored      A5
│   ├── hint/                       H1–H4
│   │   ├── standalone-pass         H1
│   │   ├── standalone-fail         H2
│   │   ├── with-literal-pass       H3
│   │   └── path-mismatch           H4 — error mentions hint:path
│   ├── contains/                   C1–C6
│   │   ├── order-free-pass         C1
│   │   ├── missing-fragment        C2
│   │   ├── full-line-required      C3
│   │   ├── start-with-prefix       C4
│   │   ├── end-with-suffix         C5
│   │   └── meta-lines-ignored      C6
│   ├── ansi-color/                 AC1–AC6, AC8–AC10 (AC7 is parse error)
│   │   ├── named-gray-pass         AC1
│   │   ├── plain-text-fails        AC2
│   │   ├── named-green-pass        AC3
│   │   ├── raw-equivalent-pass     AC4 — #90 equals gray
│   │   ├── raw-256-pass            AC5
│   │   ├── wrong-sgr-fails         AC6
│   │   ├── bold-pass               AC8
│   │   ├── bold-gray-combo         AC9
│   │   └── missing-bold-fails      AC10
│   ├── regex/                      R1–R3
│   │   ├── block-dots-pass         R1
│   │   ├── block-glued-fail        R2
│   │   └── optional-wrapper        R3 — dots absent, summary only
│   └── inline-any-of/              X1–X3
│       ├── first-branch            X1
│       ├── second-branch           X2
│       └── prefix-required         X3
├── modes/                          match options & normalization (N1–N3)
│   ├── match-contains/
│   │   └── embedded-in-log         N1 — contiguous subregion in long log
│   └── normalization/
│       └── crlf-normalized         N2 — \r\n actual vs \n template
└── integration/                    realistic CLI templates (§10.8)
    ├── dot-progress-summary        R1 — regex dots + summary line
    └── help-contains               R4 — help keywords via contains block
```

## Test Index

| ID | Leaf | Description |
|----|------|-------------|
| P1 | `parse/valid/literal-lines` | Two literal lines |
| P2 | `parse/valid/block-optional` | Block optional between literals |
| P3 | `parse/valid/hint-standalone` | Hint segment label `id` |
| P4 | `parse/valid/hint-with-literal` | Literal prefix + hint |
| P5 | `parse/errors/bare-tag` | Bare `<id>` parse error |
| P6 | `parse/errors/unclosed-hint` | Unclosed hint parse error |
| P7 | `parse/valid/inline-optional` | Inline optional segment |
| P8 | `parse/valid/block-any-of` | Two-branch AnyOfBlock |
| P9 | `parse/valid/contains-block` | Three-fragment ContainsBlock |
| P10 | `parse/errors/hint-label-mismatch` | Hint label mismatch |
| P11 | `parse/valid/ansi-color-named` | Named gray AnsiColor |
| P12 | `parse/valid/contains-start-with` | StartWith fragment mode |
| P13 | `parse/errors/empty-ansi-color` | Empty ansi-color inner |
| P14 | `parse/valid/ansi-color-raw` | Raw `#38;5;208` SGR |
| P15 | `parse/errors/unknown-ansi-name` | Unknown `orange` color name |
| P16 | `parse/valid/regex-block-inline` | Block + inline regex |
| P17 | `parse/valid/inline-any-of` | Inline any-of pattern line |
| P18 | `parse/valid/ansi-color-combined` | bold + gray tokens |
| P19 | `parse/valid/escaped-literal` | Escaped `<optional>` literal |
| M1 | `match/literal/exact-pass` | Exact single-line match |
| M2 | `match/literal/strict-trailing-fail` | Strict trailing newline fail |
| M2b | `match/literal/strict-trailing-pass` | Trailing newline agreement |
| M3 | `match/literal/line-mismatch` | Line-level mismatch |
| O1–O9 | `match/optional/*` | Block/inline optional semantics |
| A1–A5 | `match/any-of/*` | Branch selection and reporting |
| H1–H4 | `match/hint/*` | Literal hint matching |
| C1–C7 | `match/contains/*` | Order-free fragments and inline pattern fragments |
| AC1–AC10 | `match/ansi-color/*` + `parse/errors/empty-ansi-color` | ANSI assertions |
| R1–R3 | `match/regex/*` | Regex block + optional wrapper |
| X1–X3 | `match/inline-any-of/*` | Inline any-of scoping |
| N1 | `modes/match-contains/embedded-in-log` | MatchContains contiguous |
| N2 | `modes/normalization/crlf-normalized` | CRLF normalization |
| R1 | `integration/dot-progress-summary` | Dot progress + summary |
| R4 | `integration/help-contains` | Help text contains block |

## How to Run

```sh
doctest vet ./assert/tests/output-assert
doctest test ./assert/tests/output-assert/...
go test ./assert/...
```

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/assert"
)

type Request struct {
	Operation	string
	Template	string
	Actual		string
	Options		[]string
	ExpectParseError	bool
	ExpectMatchError	bool
}

type Response struct {
	ParseOK		bool
	ParseError	string
	ParsedSummary	string
	MatchOK		bool
	MatchError	string
}

func Run(t *testing.T, req *Request) (*Response, error) {
	resp := &Response{}

	p, parseErr := assert.Parse(req.Template)
	if parseErr != nil {
		resp.ParseError = parseErr.Error()
		if req.Operation == "parse" {
			return resp, nil
		}
		return resp, parseErr
	}
	resp.ParseOK = true
	resp.ParsedSummary = summarizePattern(p)

	if req.Operation == "parse" {
		return resp, nil
	}

	matchErr := assert.Match(p, req.Actual, buildMatchOptions(req.Options)...)
	if matchErr != nil {
		resp.MatchError = matchErr.Error()
		return resp, nil
	}
	resp.MatchOK = true
	return resp, nil
}

func buildMatchOptions(opts []string) []assert.Option {
	var out []assert.Option
	for _, o := range opts {
		switch o {
		case "contains":
			out = append(out, assert.Contains())
		case "normalize_newlines:false":
			out = append(out, assert.NormalizeNewlines(false))
		}
	}
	return out
}
```
