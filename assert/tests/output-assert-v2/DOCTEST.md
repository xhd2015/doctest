# Output Assert DSL v2 Tests

## Version
0.0.2

## DSN (Domain Specific Notion)

### Participants

- **Author** — writes a v2 template with YAML frontmatter (`version: 2`, placeholder
  defs) and a strict line-by-line body (`__PLACEHOLDER__`, `...N lines omitted...`,
  `<ansi-color>` spans, regex-intent lines).
- **Facade** — `assert.Parse` / `assert.Match` auto-detect v2 via `version: 2` in
  YAML header; otherwise dispatch to `legacy_v1` unchanged.
- **v2 Parser** — reads YAML header + body, builds immutable `Pattern` AST. Rejects
  undefined placeholders, unknown types, bad omit syntax, invalid YAML.
- **v2 Matcher** — strict sequential line-by-line match: pattern lines, regex lines,
  omit markers, color spans, typed placeholders. No `MatchContains` in v2.
- **legacy_v1 Parser/Matcher** — unchanged v1 tags (`<contains>`, `<ansi-color>`, …)
  for templates without a v2 header.

### Behaviors

- **Version detect** — leading `---` + `version: 2` → v2 path; else → legacy_v1.
- **Parse** — template text → `Pattern` or parse error with message.
- **Match** — `Pattern` + actual string → nil (pass) or rich diff error (line numbers).
- **Placeholder expand** — `type=string` → `[^\n]*`; `type=number` → `-?\d+(\.\d+)?`.
- **Regex detect** — protected-region scan + strong-signal intent (not naive metachar grep).
- **Omit** — `...N lines omitted...` consumes exactly N actual lines (N≥0).
- **Color** — `<ansi-color SPEC>text</ansi-color>` asserts strict ANSI open + `\x1b[0m` reset.

## Decision Tree

```
output-assert-v2
├── parse/                          Operation = parse
│   ├── valid/                      successful v2 AST shapes (V2-P1–P6)
│   │   ├── placeholder-number      V2-P1 — __PORT__: type=number
│   │   ├── placeholder-string      V2-P2 — __USER__: type=string
│   │   ├── compact-yaml-form       V2-P3 — k=v metadata + human explanation
│   │   ├── regex-line              V2-P4 — regex-intent line classified
│   │   ├── omit-marker             V2-P5 — ...3 lines omitted...
│   │   ├── literal-only            V2-P6 — version-only header + literals
│   │   └── leading-blank-before-header V2-P9 — trim blank lines before ---
│   └── errors/                     parse must fail or v1 fallback (V2-E1–E4)
│       ├── undefined-placeholder   V2-E1 — __MISSING__ not in header
│       ├── unknown-type            V2-E2 — type=boolean rejected
│       ├── bad-omit-syntax         V2-E3 — non-numeric omit count
│       └── missing-version         V2-E4 — no version:2 → legacy_v1 literal
├── match/                          Operation = match (strict full match only)
│   ├── placeholders/               typed __NAME__ expansion (V2-M1–M3, M14)
│   │   ├── number-pass             V2-M1
│   │   ├── number-fail             V2-M2
│   │   ├── string-pass             V2-M3
│   │   └── number-float            V2-M14 — float placeholder
│   ├── regex/                      regex lines (V2-M4–M5, M13)
│   │   ├── line-pass               V2-M4
│   │   ├── line-fail               V2-M5
│   │   └── alternation             V2-M13 — (ok|fail)
│   ├── omit/                       omit marker consumption (V2-M6–M7)
│   │   ├── three-lines             V2-M6
│   │   └── wrong-count             V2-M7
│   ├── strict/                     strict sequential policy (V2-M8–M9)
│   │   ├── extra-line              V2-M8
│   │   ├── trailing-newline        V2-M9
│   │   └── trailing-blank-lines    V2-M17 — trailing blank body lines preserved
│   ├── literal-preservation/       regex-detection edge cases (V2-M10–M12)
│   │   ├── with-dot                V2-M10 — version 1.0 stays pattern
│   │   ├── dollar                  V2-M11 — cost: $5.00 mid-line $
│   │   └── parens                  V2-M12 — (1 Cached) stays pattern
│   └── color/                      v2 <ansi-color> spans (V2-M15–M16)
│       ├── gray-pass               V2-M15
│       └── plain-fail              V2-M16
└── integration/                    realistic + compatibility (V2-I1–I2 + 188 CLI cookbook)
    ├── server-startup-log          V2-I1 — PORT + omit + color status
    ├── v1-still-works              V2-I2 — v1 <contains> via legacy_v1
    └── real-world/                 188 simulated CLI transcripts (17 categories)
        ├── unix-text/              cat, grep, rg, head, tail, sed, awk, find, …
        ├── go-toolchain/           build, test, mod, vet, fmt, get, install, …
        ├── rust-toolchain/         cargo, rustc, rustup, clippy, bench, …
        ├── node-js/                npm, yarn, pnpm, eslint, jest, vite, tsc, …
        ├── python/                 pip, pytest, ruff, black, mypy, poetry, uv, …
        ├── git/                    status, log, diff, clone, merge, stash, …
        ├── http-clients/           curl, wget, httpie
        ├── containers/             docker, compose, kubectl, helm, podman
        ├── build-systems/          make, cmake, ninja, bazel, gradle, maven, …
        ├── databases/              psql, mysql, redis, sqlite3, mongosh, pg_dump
        ├── jvm-kotlin/             javac, kotlin, gradle, mvn, sbt, scala
        ├── c-cpp/                  gcc, g++, clang, lldb, gdb, valgrind
        ├── shell/                  bash, sh, zsh, fish, dash
        ├── package-managers/       brew, apt, pacman, apk, dnf, nix, mise
        ├── cloud-infra/            terraform, aws, gcloud, az, pulumi, gh, glab
        ├── languages-other/        ruby, php, swift, deno, bun, crystal, zig
        └── misc-devtools/          openssl, ssh, jq, yq, helm lint, ffmpeg, protoc
```

## Test Index

| ID | Leaf | Description |
|----|------|-------------|
| V2-P1 | `parse/valid/placeholder-number` | Number placeholder in header + body |
| V2-P2 | `parse/valid/placeholder-string` | String placeholder |
| V2-P3 | `parse/valid/compact-yaml-form` | Compact k=v metadata + explanation |
| V2-P4 | `parse/valid/regex-line` | Regex-intent line in AST |
| V2-P5 | `parse/valid/omit-marker` | Omit marker line |
| V2-P6 | `parse/valid/literal-only` | Literals only under v2 header |
| V2-P9 | `parse/valid/leading-blank-before-header` | Leading blank lines before `---` trimmed |
| V2-E1 | `parse/errors/undefined-placeholder` | Undefined __MISSING__ |
| V2-E2 | `parse/errors/unknown-type` | Unknown placeholder type |
| V2-E3 | `parse/errors/bad-omit-syntax` | Invalid omit count |
| V2-E4 | `parse/errors/missing-version` | No version:2 → v1 literal fallback |
| V2-M1 | `match/placeholders/number-pass` | PORT=8901 matches |
| V2-M2 | `match/placeholders/number-fail` | non-number PORT fails |
| V2-M3 | `match/placeholders/string-pass` | USER string matches |
| V2-M4 | `match/regex/line-pass` | .*middle.*suffix regex pass |
| V2-M5 | `match/regex/line-fail` | regex line mismatch |
| V2-M6 | `match/omit/three-lines` | omit consumes 3 middle lines |
| V2-M7 | `match/omit/wrong-count` | omit count mismatch |
| V2-M8 | `match/strict/extra-line` | extra actual line rejected |
| V2-M9 | `match/strict/trailing-newline` | strict trailing newline |
| V2-M17 | `match/strict/trailing-blank-lines` | trailing blank body lines preserved |
| V2-M10 | `match/literal-preservation/with-dot` | version 1.0 literal |
| V2-M11 | `match/literal-preservation/dollar` | cost: $5.00 literal |
| V2-M12 | `match/literal-preservation/parens` | (1 Cached) literal |
| V2-M13 | `match/regex/alternation` | (ok\|fail) alternation |
| V2-M14 | `match/placeholders/number-float` | float MS placeholder |
| V2-M15 | `match/color/gray-pass` | <ansi-color gray> pass |
| V2-M16 | `match/color/plain-fail` | missing ANSI fails |
| V2-I1 | `integration/server-startup-log` | PORT + omit stack + status |
| V2-I2 | `integration/v1-still-works` | v1 contains via legacy_v1 |
| V2-RW | `integration/real-world/*` | **188** CLI cookbook leaves — see category dirs above |

Regenerate cookbook leaves:

```sh
go run ./script/generate/real-world-assert-cases/main.go
```

## How to Run

```sh
doctest vet ./assert/tests/output-assert-v2
doctest test ./assert/tests/output-assert-v2/...
go test ./assert/...
doctest test ./assert/tests/output-assert/...   # no v1 regressions
```

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/assert"
)

type Request struct {
	Operation         string
	Template          string
	Actual            string
	ExpectParseError  bool
	ExpectV1Fallback  bool
}

type Response struct {
	ParseOK         bool
	ParseError      string
	ParsedSummary   string
	MatchOK         bool
	MatchError      string
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