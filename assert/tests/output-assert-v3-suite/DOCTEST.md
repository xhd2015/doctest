# Output Assert DSL v3 Full Suite Tests

## Version
0.0.2

## DSN (Domain Specific Notion)

### Participants

- **Author** — writes a v3 template with YAML frontmatter (`version: 3` or
  dialect header without version key; placeholder defs `type=` / `regex=`) and a
  strict line-by-line body. Content lines are **raw Go regex**; omit markers and
  `<ansi-color>` spans remain special. Literal dots/parens/dollars must be escaped
  (`0\.001s`, `\(1 Cached\)`, `\$5\.00`). Color-tag **inner** text stays literal
  (engine applies QuoteMeta).
- **Facade** — `assert.Parse` / `assert.Match` route templates:
  - explicit `version: 2` → `legacy_v2`
  - v1 tag DSL (no YAML version-2/3 dialect) → `legacy_v1`
  - explicit `version: 3` **or** YAML dialect with placeholders and no version → **v3**
  - unknown version value → parse error
- **v3 Parser** — reads YAML header + body, builds immutable v3 Pattern
  (placeholders + items). Content lines become `RegexLine` (`^…$`). Placeholders
  inject named capture groups. Rejects dual `type=`+`regex=`, invalid `regex=`
  fragments, undefined placeholders, unknown types, bad omit syntax, unknown version.
- **v3 Matcher** — strict sequential match; same-value bindings for repeated
  `__NAME__`; omit consumes exactly N lines; color spans keep structure with
  QuoteMeta on inner text.
- **legacy_v2 / legacy_v1** — still available for routing/smoke; this suite is
  v3-primary (no production leaf requires `version: 2`).

### Behaviors

- **Version detect** — `version: 3` or placeholder dialect without version → v3;
  `version: 2` → legacy_v2; non-dialect / pure v1 tags → legacy_v1.
- **Content lines** — raw Go regex full-line match. Author escapes metachars for
  literals. No `hasRegexIntent` classifier in v3.
- **Placeholders** — `type=string` → non-greedy `[^\n]*?`; `type=number` →
  `-?\d+(?:\.\d+)?`; `regex=<fragment>` → custom; dual type+regex → parse error.
- **Same-value binding** — repeated `__NAME__` must capture the same string.
- **Omit** — `...N lines omitted...` is special (not content regex).
- **Color** — `<ansi-color SPEC>inner</ansi-color>`: structure kept; inner
  QuoteMeta'd; outside tags remain raw RE.
- **Parse / Match** — template → Pattern or parse error; Pattern + actual → nil or error.

## Mapping notes (P3: output-assert-v2 → v3)

| Before (v2 suite) | After (this suite) |
|-------------------|--------------------|
| Tree path `assert/tests/output-assert-v2` | Renamed to `assert/tests/output-assert-v3-suite` (canonical full suite + cookbook) |
| Helper `v2Template` → `version: 2` | Helper `v3Template` → `version: 3` |
| Pattern lines (literal via hasRegexIntent) | Escaped raw RE (`\.`, `\(`, `\$`, …) |
| Intentional regex leaves (`.*…`, `(ok\|fail)`) | Unchanged raw RE bodies under `match/regex/`, `parse/valid/regex-line` |
| Color inner dots | Still unescaped (engine QuoteMeta) |
| Summary `LiteralLine` for content | Summary `RegexLine` / `RegexLine+…` |
| `missing-version` = no `version:2` → v1 | Unchanged semantics: non-placeholder YAML frontmatter → legacy_v1 |
| Focused unit suite | Remains at `assert/tests/output-assert-v3` (binding, raw-dot, routing smoke) |
| v1 tree | `assert/tests/output-assert` untouched |

Companion suite `output-assert-v3` holds engine-edge leaves (binding, dual type+regex,
`version-2-still-works` legacy smoke). This suite holds the migrated cookbook and
broader match/parse coverage under v3 defaults.

## Decision Tree

```
output-assert-v3-suite          (was output-assert-v2; v3 templates)
├── parse/                          Operation = parse
│   ├── valid/                      successful v3 AST shapes (V3S-P1–P9)
│   │   ├── placeholder-number      V3S-P1 — __PORT__: type=number
│   │   ├── placeholder-string      V3S-P2 — __USER__: type=string
│   │   ├── compact-yaml-form       V3S-P3 — k=v metadata + human explanation
│   │   ├── regex-line              V3S-P4 — raw RE content line
│   │   ├── omit-marker             V3S-P5 — ...3 lines omitted...
│   │   ├── literal-only            V3S-P6 — version-only header + escaped/plain lines
│   │   └── leading-blank-before-header V3S-P9 — trim blank lines before ---
│   └── errors/                     parse must fail or v1 fallback (V3S-E1–E4)
│       ├── undefined-placeholder   V3S-E1 — __MISSING__ not in header
│       ├── unknown-type            V3S-E2 — type=boolean rejected
│       ├── bad-omit-syntax         V3S-E3 — non-numeric omit count
│       └── missing-version         V3S-E4 — non-dialect YAML → legacy_v1
├── match/                          Operation = match (strict full match only)
│   ├── placeholders/               typed __NAME__ expansion (V3S-M1–M3, M14)
│   │   ├── number-pass             V3S-M1
│   │   ├── number-fail             V3S-M2
│   │   ├── string-pass             V3S-M3
│   │   └── number-float            V3S-M14 — float placeholder
│   ├── regex/                      intentional raw RE (V3S-M4–M5, M13)
│   │   ├── line-pass               V3S-M4
│   │   ├── line-fail               V3S-M5
│   │   └── alternation             V3S-M13 — (ok|fail)
│   ├── omit/                       omit marker consumption (V3S-M6–M7)
│   │   ├── three-lines             V3S-M6
│   │   └── wrong-count             V3S-M7
│   ├── strict/                     strict sequential policy (V3S-M8–M9, M17)
│   │   ├── extra-line              V3S-M8
│   │   ├── trailing-newline        V3S-M9
│   │   └── trailing-blank-lines    V3S-M17
│   ├── literal-preservation/       escaped literals under v3 raw RE (V3S-M10–M12)
│   │   ├── with-dot                V3S-M10 — version 1\.0
│   │   ├── dollar                  V3S-M11 — cost: \$5\.00
│   │   └── parens                  V3S-M12 — \(1 Cached\)
│   └── color/                      <ansi-color> spans (V3S-M15–M16)
│       ├── gray-pass               V3S-M15 — inner unescaped
│       └── plain-fail              V3S-M16
└── integration/                    realistic + compatibility (V3S-I1–I2 + 188 CLI cookbook)
    ├── server-startup-log          V3S-I1 — PORT + omit + color status
    ├── v1-still-works              V3S-I2 — v1 <contains> via legacy_v1
    └── real-world/                 188 simulated CLI transcripts (17 categories)
        ├── unix-text/ …
        ├── go-toolchain/ …
        ├── rust-toolchain/ …
        ├── node-js/ …
        ├── python/ …
        ├── git/ …
        ├── http-clients/ …
        ├── containers/ …
        ├── build-systems/ …
        ├── databases/ …
        ├── jvm-kotlin/ …
        ├── c-cpp/ …
        ├── shell/ …
        ├── package-managers/ …
        ├── cloud-infra/ …
        ├── languages-other/ …
        └── misc-devtools/ …
```

## Test Index

| ID | Leaf | Description |
|----|------|-------------|
| V3S-P1 | `parse/valid/placeholder-number` | Number placeholder in header + body |
| V3S-P2 | `parse/valid/placeholder-string` | String placeholder |
| V3S-P3 | `parse/valid/compact-yaml-form` | Compact k=v metadata + explanation |
| V3S-P4 | `parse/valid/regex-line` | Raw RE content line in AST |
| V3S-P5 | `parse/valid/omit-marker` | Omit marker line |
| V3S-P6 | `parse/valid/literal-only` | Two content lines → RegexLine+RegexLine |
| V3S-P9 | `parse/valid/leading-blank-before-header` | Leading blank lines before `---` trimmed |
| V3S-E1 | `parse/errors/undefined-placeholder` | Undefined __MISSING__ |
| V3S-E2 | `parse/errors/unknown-type` | Unknown placeholder type |
| V3S-E3 | `parse/errors/bad-omit-syntax` | Invalid omit count |
| V3S-E4 | `parse/errors/missing-version` | Non-dialect YAML → v1 literal fallback |
| V3S-M1 | `match/placeholders/number-pass` | PORT=8901 matches |
| V3S-M2 | `match/placeholders/number-fail` | non-number PORT fails |
| V3S-M3 | `match/placeholders/string-pass` | USER string matches |
| V3S-M4 | `match/regex/line-pass` | .*middle.*suffix regex pass |
| V3S-M5 | `match/regex/line-fail` | regex line mismatch |
| V3S-M6 | `match/omit/three-lines` | omit consumes 3 middle lines |
| V3S-M7 | `match/omit/wrong-count` | omit count mismatch |
| V3S-M8 | `match/strict/extra-line` | extra actual line rejected |
| V3S-M9 | `match/strict/trailing-newline` | strict trailing newline |
| V3S-M17 | `match/strict/trailing-blank-lines` | trailing blank body lines preserved |
| V3S-M10 | `match/literal-preservation/with-dot` | escaped `version 1\.0` |
| V3S-M11 | `match/literal-preservation/dollar` | escaped `cost: \$5\.00` |
| V3S-M12 | `match/literal-preservation/parens` | escaped `\(1 Cached\)` |
| V3S-M13 | `match/regex/alternation` | (ok\|fail) alternation |
| V3S-M14 | `match/placeholders/number-float` | float MS placeholder |
| V3S-M15 | `match/color/gray-pass` | <ansi-color gray> pass |
| V3S-M16 | `match/color/plain-fail` | missing ANSI fails |
| V3S-I1 | `integration/server-startup-log` | PORT + omit stack + status |
| V3S-I2 | `integration/v1-still-works` | v1 contains via legacy_v1 |
| V3S-RW | `integration/real-world/*` | **188** CLI cookbook leaves — v3 escaped templates |

## How to Run

Leaves that need skip-from-default discovery should use `label: e2e` (public L3) so `./...` discovery **skips** this
cookbook (keeps cold full-module self-tests under ~3 minutes). Opt in explicitly:

```sh
doctest vet ./assert/tests/output-assert-v3-suite
doctest test ./assert/tests/output-assert-v3-suite --label e2e
doctest test ./assert/tests/output-assert-v3          # focused engine suite (not heavy)
doctest test ./assert/tests/output-assert            # v1 unchanged
go test ./assert/...
```

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
	ExpectV1Fallback bool
}

type Response struct {
	ParseOK       bool
	ParseError    string
	ParsedSummary string
	MatchOK       bool
	MatchError    string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
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
