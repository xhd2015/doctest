# Output Assert DSL v3 Tests

## Version
0.0.2

## DSN (Domain Specific Notion)

### Participants

- **Author** — writes a v3 template with YAML frontmatter (`version: 3` or
  dialect header without version key; placeholder defs `type=` / `regex=`) and a
  strict line-by-line body. Content lines are raw Go regex; omit markers and
  `<ansi-color>` spans remain special.
- **Facade** — `assert.Parse` / `assert.Match` route templates:
  - explicit `version: 2` → `legacy_v2`
  - v1 tag DSL (no YAML version-2/3 dialect) → `legacy_v1`
  - explicit `version: 3` **or** YAML dialect header without version → **v3**
  - unknown version value → parse error
- **v3 Parser** — reads YAML header + body, builds immutable v3 Pattern
  (`v3Pattern` with placeholders + items). Content lines become `RegexLine`
  (`^…$`). Placeholders inject named capture groups. Rejects dual `type=`+`regex=`,
  invalid `regex=` fragments, undefined placeholders, unknown types, bad omit
  syntax, unknown version.
- **v3 Matcher** — strict sequential match; tracks same-value bindings for
  repeated `__NAME__`; omit consumes exactly N lines; color spans keep structure
  with QuoteMeta on inner text (literal dots etc.).
- **legacy_v2 / legacy_v1** — unchanged engines for `version: 2` and pure v1 tags.

### Behaviors

- **Version detect** — `version: 3` or YAML dialect without version → v3;
  `version: 2` → legacy_v2; no YAML dialect → legacy_v1; unknown version → error.
- **Content lines** — raw Go regex full-line match (`^…$`). Author escapes
  metachars (`0.001s` → `0\.001s`). No `hasRegexIntent` in v3.
- **Placeholders** — `type=string` → `[^\n]*?`; `type=number` →
  `-?\d+(?:\.\d+)?`; `regex=<fragment>` → custom subpattern; both type+regex →
  parse error.
- **Same-value binding** — repeated `__NAME__` must capture the same string.
- **Omit** — `...N lines omitted...` is special (not content regex).
- **Color** — `<ansi-color SPEC>inner</ansi-color>`: structure kept; inner text
  QuoteMeta'd; outside tags remain raw RE.
- **Parse** — template text → Pattern or parse error.
- **Match** — Pattern + actual → nil (pass) or match error (bindings, counts,
  line mismatch).

## Decision Tree

```
output-assert-v3
├── parse/                              Operation = parse
│   ├── valid/                          successful v3 AST / dialect shapes
│   │   ├── version-3-header            explicit version: 3
│   │   ├── default-no-version          YAML placeholders, no version key → v3
│   │   ├── placeholder-string          __USER__: type=string
│   │   ├── placeholder-number          __PORT__: type=number
│   │   ├── placeholder-regex-custom    __ID__: regex=[a-z]+
│   │   ├── omit-marker                 ...3 lines omitted...
│   │   └── compact-yaml-form           k=v metadata + human explanation
│   └── errors/                         parse must fail
│       ├── dual-type-and-regex         type= and regex= on same placeholder
│       ├── invalid-regex-fragment      regex=[ does not compile
│       ├── undefined-placeholder       __MISSING__ not in header
│       ├── unknown-type                type=boolean without regex=
│       ├── bad-omit-syntax             non-numeric omit count
│       └── unknown-version             version: 9 rejected
├── match/                              Operation = match (strict full match)
│   ├── raw-regex/                      content lines are raw RE (no hasRegexIntent)
│   │   ├── raw-dot-is-any              template a.c matches actual aXc
│   │   └── escaped-dot-literal         a\.c matches only a.c
│   ├── placeholders/                   typed expansion + non-greedy string
│   │   ├── number-loose                1, -2, 3.14 via type=number
│   │   └── string-midline              non-greedy mid-line with trailing content
│   ├── binding/                        same-value binding across repeats
│   │   ├── same-ok                     two __ID__ same value
│   │   └── same-fail                   two __ID__ different → error names ID
│   ├── omit/                           omit marker consumption
│   │   ├── three-lines                 omit consumes 3 middle lines
│   │   └── wrong-count                 omit count mismatch
│   ├── color/                          ansi-color + placeholder composition
│   │   └── ansi-with-placeholder       inner color text literal (dots)
│   └── strict/                         strict sequential policy
│       └── extra-line                  extra actual line rejected
└── routing/                            version dispatch compatibility
    ├── version-2-still-works           version: 2 → legacy_v2 still matches
    └── v1-still-works                  pure v1 tags → legacy_v1
```

## Test Index

| ID | Leaf | Description |
|----|------|-------------|
| V3-P1 | `parse/valid/version-3-header` | Explicit `version: 3` routes to v3 |
| V3-P2 | `parse/valid/default-no-version` | YAML dialect without version → v3 |
| V3-P3 | `parse/valid/placeholder-string` | String placeholder in header + body |
| V3-P4 | `parse/valid/placeholder-number` | Number placeholder |
| V3-P5 | `parse/valid/placeholder-regex-custom` | `regex=` custom subpattern |
| V3-P6 | `parse/valid/omit-marker` | Omit marker line in AST |
| V3-P7 | `parse/valid/compact-yaml-form` | Compact k=v metadata + explanation |
| V3-E1 | `parse/errors/dual-type-and-regex` | Both type= and regex= → parse error |
| V3-E2 | `parse/errors/invalid-regex-fragment` | Invalid regex= fragment |
| V3-E3 | `parse/errors/undefined-placeholder` | Undefined `__MISSING__` |
| V3-E4 | `parse/errors/unknown-type` | Unknown type without regex |
| V3-E5 | `parse/errors/bad-omit-syntax` | Invalid omit count |
| V3-E6 | `parse/errors/unknown-version` | `version: 9` → parse error |
| V3-M1 | `match/raw-regex/raw-dot-is-any` | Unescaped `.` is any-char RE |
| V3-M2 | `match/raw-regex/escaped-dot-literal` | `\.` is literal dot |
| V3-M3 | `match/placeholders/number-loose` | type=number matches 1, -2, 3.14 |
| V3-M4 | `match/placeholders/string-midline` | Non-greedy string + trailing content |
| V3-M5 | `match/binding/same-ok` | Repeated `__ID__` same value |
| V3-M6 | `match/binding/same-fail` | Different values → match error names ID |
| V3-M7 | `match/omit/three-lines` | omit consumes 3 middle lines |
| V3-M8 | `match/omit/wrong-count` | omit count mismatch |
| V3-M9 | `match/color/ansi-with-placeholder` | Color inner QuoteMeta + placeholder |
| V3-M10 | `match/strict/extra-line` | Extra actual line rejected |
| V3-R1 | `routing/version-2-still-works` | version: 2 still uses legacy_v2 |
| V3-R2 | `routing/v1-still-works` | v1 `<contains>` via legacy_v1 |

## How to Run

```sh
doctest vet ./assert/tests/output-assert-v3
doctest test ./assert/tests/output-assert-v3
go test ./assert/...
```

Classic TDD: new v3 leaves are expected **RED** until the v3 engine is implemented.
Routing leaves (`version-2-still-works`, `v1-still-works`) may already be GREEN.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/assert"
)

type Request struct {
	Operation        string
	Template         string
	Actual           string
	ExpectParseError bool
}

type Response struct {
	ParseOK       bool
	ParseError    string
	ParsedSummary string
	MatchOK       bool
	MatchError    string
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
	resp.ParsedSummary = assert.PatternSummary(p)

	if req.Operation == "parse" {
		return resp, nil
	}

	matchErr := assert.Match(p, req.Actual)
	if matchErr != nil {
		resp.MatchError = matchErr.Error()
		return resp, nil
	}
	resp.MatchOK = true
	return resp, nil
}
```
