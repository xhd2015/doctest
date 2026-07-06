---
name: doctest-output-assert
description: >-
  Output assert DSL v2 for doctest ASSERT.md — YAML header, placeholders, strict
  line match, and migration from strings.Contains. Use when asserting CLI or text
  output in doc-style tests.
---

# Output Assert DSL

Package: `github.com/xhd2015/doctest/assert`

Readable **template-shaped expected output** for doctest `ASSERT.md` leaves. Replaces ad-hoc `strings.Contains`, `strings.Index`, and hand-rolled parsing.

Design history: `doc/DESIGN_OUTPUT_ASSERT.md`. v1 tag syntax lives in `assert/legacy_v1/`.

## When to use

| Situation | Approach |
|-----------|----------|
| Bounded stdout/stderr (typical leaf) | `assert.Output(t, actual, template)` — strict full match |
| Variable port, path, user, count | `__PLACEHOLDER__` in YAML header + template body |
| Middle section may differ (stack trace, progress) | `...N lines omitted...` |
| Flexible single line (dots, prefixes) | regex line (e.g. `^\.+$`) |
| ANSI-colored segments | `<ansi-color bold gray>…</ansi-color>` |
| Document example values for readers | `example=…` in placeholder header (metadata only) |

**v2 is strict line-by-line full match.** There is no `<contains>` or `assert.Contains()` in v2 templates.

## API

```go
import "github.com/xhd2015/doctest/assert"

p, err := assert.Parse(template)
err = assert.Match(p, actual)       // strict full match (default)
assert.Output(t, actual, template)  // Parse + Match + t.Fatal
```

- Templates with `version: 2` in the YAML header use the v2 parser.
- Templates without `version: 2` fall back to v1 (`assert/legacy_v1/`).
- `\r\n` → `\n` normalization is always on.
- Trailing newlines are **strict** (template and actual must agree).
- Matching is **case-sensitive**.

### CLI stdout trailing newline (required)

User-facing CLI stdout **must** end with a final `\n` (POSIX convention — keeps
the shell prompt on its own line). v2 templates for CLI output **must** include
that trailing newline too, so template and actual agree **and** the product
behaves correctly in a real terminal.

**Go authoring pattern** — put the closing backtick on the line **after** the
last content line (not on the same line, and not separated by an extra blank
line):

````go
assert.Output(t, resp.Stdout, `---
version: 2
---
Hello world
`)
````

The newline before `` ` `` makes the raw string end with `\n` without adding an
extra blank output line. A **blank line** between the last content line and the
backtick adds an empty line to the template body — do not do that unless the
product output truly ends with a blank line.

**Trap to avoid:** `` `...last line` `` on one line omits trailing `\n` from the
template. That forces the implementation to also omit `\n`, which glues the
shell prompt to the last line. Implementers must not strip `\n` from product code
to pass such tests — fix the template instead.

## Template shape (v2)

```DSL
---
version: 2
__PORT__: type=number, example=8901, a port
__USER__: type=string, example=alice, logged-in user
---
Server listen on: __PORT__
...3 lines omitted...
Hello __USER__
<ansi-color bold gray>1 Cached</ansi-color>
```

| Part | Role |
|------|------|
| `---` … `---` | YAML header; `version: 2` selects v2 |
| `__NAME__:` | Placeholder definition (compact `k=v` metadata + human explanation) |
| Body lines | Strict sequential match against actual output |

In Go `ASSERT.md` blocks, pass the template as one raw string starting with `---` — no `` + `` concatenation needed. Leading blank lines before the opening `---` are trimmed; trailing blank lines in the body are preserved.

### Placeholder definitions

Compact form (preferred):

```yaml
__PORT__: type=number, example=8901, a port
```

- `k=v` pairs before the first bare word — machine metadata
- Trailing text after metadata — human explanation only (ignored by matcher)

| Field | Required | Values |
|-------|----------|--------|
| `type` | yes | `string`, `number` |
| `example` | no | documentation for readers |

| type | Matches on one line |
|------|---------------------|
| `string` | any non-newline text |
| `number` | integer or float (`-?\d+(\.\d+)?`) |

Every `__NAME__` used in the body must be defined in the header.

### Body line kinds

1. **Pattern line** (default) — literal text with optional `__PLACEHOLDER__` and `<ansi-color>` spans. Full line must match.

2. **Regex line** — detected by regex-intent scan (see below). Entire line is a Go regexp (full line match). Placeholders expand to typed subpatterns; `<ansi-color>` spans expand to literal ANSI envelope.

3. **Omit marker** — `...N lines omitted...` (whitespace allowed around parts). Skips exactly **N** actual lines (any content). `N` must be a non-negative integer.

### Regex line detection

Before scanning, mask protected regions (`__PLACEHOLDER__`, `<ansi-color>…</ansi-color>`). A line is a **regex line** when the remaining text has a **strong regex-intent signal**:

| Signal | Examples |
|--------|----------|
| Dot-quantifier | `.*`, `.+`, `.?` |
| Anchors | `^` at start; `$` at end only |
| Escape atoms | `\d`, `\w`, `\s`, `\b`, … |
| Char class | `[a-z]+` |
| Alternation | `(ok\|fail)`, `foo\|bar` |
| Braced quantifier | `a{2,3}` |

**Stays literal** without a strong signal: `version 1.0`, `file.go:42`, `cost: $5.00`, `(1 Cached)`.

Examples:

```DSL
^\.+$                          # regex — progress dots line
.*Some middle content.*suffix  # regex — flexible middle
version 1.0                    # pattern — version dot is literal
```

### Color spans (only inline tag in v2)

`<ansi-color SPEC>inner text</ansi-color>` — same tag and token set as v1.

Space-separated tokens; **strict** open SGR + `\x1b[0m` reset immediately after inner text.

| Token | Open SGR | Example |
|-------|----------|---------|
| `bold` | `\x1b[1m` | `<ansi-color bold>title</ansi-color>` |
| `red` | `\x1b[31m` | `<ansi-color red>FAIL</ansi-color>` |
| `green` | `\x1b[32m` | `<ansi-color green>OK</ansi-color>` |
| `gray` | `\x1b[90m` | `<ansi-color gray>1 Cached</ansi-color>` |
| `#` + params | `\x1b[<params>m` | `<ansi-color #38;5;208>warn</ansi-color>` |

Combined — emitted left to right: `<ansi-color bold gray>1 Cached</ansi-color>`.

## ASSERT.md pattern

**Recommended** (not required by `doctest vet`): mirror the template in prose:

````markdown
## Expected Output

```
---
version: 2
__PORT__: type=number, example=8901, a port
---
Server listen on: __PORT__
```

## Expected
- stdout reports listening port

```go
import (
    "testing"
    "github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatal(err)
    }
    assert.Output(t, resp.Stdout, `---
version: 2
__PORT__: type=number, example=8901, a port
---
Server listen on: __PORT__
`)
}
```
````

## Key semantics

### Strict sequential match

Template line 1 matches actual line 1, then line 2, etc. Omit markers consume N lines before the next template line. No extra or missing lines at the end.

### Omit markers

```DSL
Hello __USER__
...3 lines omitted...
Nice to meet you
```

The three middle actual lines may be anything (including blanks). The line count must be exact.

### Regex vs pattern

Regex applies to **one line only** — never crosses newlines. Use `...N lines omitted...` to skip variable middle sections instead of multiline regex.

## Migration from legacy asserts

| Legacy pattern | v2 prefer |
|----------------|-----------|
| `strings.Contains` loop | strict template with `...N lines omitted...` or regex line |
| `strings.Index` + `strings.Count` for dots | `^\.+$` regex line |
| Dual `Contains` for platform errors | `(linux-msg\|darwin-msg)` regex alternation |
| `metricIsColored` / `stripANSI` helpers | `<ansi-color bold gray>…</ansi-color>` |
| v1 `<hint:path>…</hint:path>` | `example=~/proj` in YAML header + literal or `__PATH__` |
| v1 `<any-of>` | regex alternation on one line |
| v1 `<optional>` | `...0 lines omitted...` or separate templates |

Existing v1 templates continue to work (no `version: 2` header). Prefer v2 for new tests.

## Real-world CLI cookbook

CLI templates in this section should end with a trailing `\n` in both the
template and the simulated actual bytes. When authoring `ASSERT.md` Go blocks,
use the closing-backtick-on-next-line pattern from **CLI stdout trailing newline**
above.

**188** doctest leaves under `assert/tests/output-assert-v2/integration/real-world/`
(17 categories). All use **simulated** bytes — no subprocess. Regenerate via
`go run ./script/generate/real-world-assert-cases/main.go`.

Categories: `unix-text`, `go-toolchain`, `rust-toolchain`, `node-js`, `python`,
`git`, `http-clients`, `containers`, `build-systems`, `databases`, `jvm-kotlin`,
`c-cpp`, `shell`, `package-managers`, `cloud-infra`, `languages-other`,
`misc-devtools`.

**YAML tip:** quote placeholder defs when `example=` contains colons:
`__LINE__: 'type=string, example=go: creating module foo'`.

Samples below; every leaf in the tree is a copy-pasteable template.

### cat — literal file dump

```DSL
---
version: 2
---
# My Project
Version 1.0
```

### grep -n — line number + match

```DSL
---
version: 2
__LINE__: type=number, example=3, 1-based line number
---
__LINE__:func main() {
```

### rg — path:line:column (no heading)

Use one placeholder for the variable prefix when literals follow on the same line
(string placeholders are line-greedy).

```DSL
---
version: 2
__HIT__: type=string, example=src/main.go:42:5, ripgrep hit prefix
---
__HIT__
```

### go build — compile error with omitted notes

```DSL
---
version: 2
__PACKAGE__: type=string, example=example.com/foo, module import path
---
# __PACKAGE__
./main.go:10:2: undefined: Bar
...2 lines omitted...
FAIL
```

### go test — pass summary with colored PASS

Keep stable substrings literal; parameterize only the trailing timing field.

```DSL
---
version: 2
__SECONDS__: type=number, example=0.123, elapsed seconds
---
=== RUN   TestFoo
--- PASS: TestFoo (0.00s)
<ansi-color green>PASS</ansi-color>
ok  	example.com/foo	__SECONDS__s
```

### go mod init

```DSL
---
version: 2
__MODULE__: type=string, example=example.com/myproject, new module path
---
go: creating new go.mod: module __MODULE__
```

### npm run build — script banner + omitted log + timing

Combine name@version into one placeholder when followed by more text on the line.

```DSL
---
version: 2
__BANNER__: type=string, example=myapp@1.0.0 build, lifecycle banner
__SECONDS__: type=number, example=1.2, build duration
---
> __BANNER__
> tsc
...4 lines omitted...
Done in __SECONDS__s.
```

### npm init — path + omitted JSON body

Placeholders absorb to end-of-line — put trailing punctuation on the next line or omit it.

```DSL
---
version: 2
__PATH__: type=string, example=/tmp/proj/package.json, output path
---
Wrote to __PATH__
...4 lines omitted...
}
```

### curl -i — HTTP headers + omitted body

```DSL
---
version: 2
__CODE__: type=number, example=200, HTTP status code
---
HTTP/1.1 __CODE__ OK
Content-Type: application/json
...3 lines omitted...
{"status":"ok"}
```

## Examples in repo

```sh
doctest test ./assert/tests/output-assert-v2/integration/real-world/...  # 188 CLI leaves
doctest test ./assert/tests/output-assert-v2/...   # 216 total v2
doctest test ./assert/tests/output-assert/...        # v1 legacy
go test ./assert/...
```