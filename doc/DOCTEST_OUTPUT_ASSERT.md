---
name: doctest-output-assert
description: >-
  Output assert DSL v3 for doctest ASSERT.md — YAML header, raw per-line regex,
  placeholders, same-value binding, and migration from ad-hoc string checks.
  Use when asserting CLI or text output in doc-style tests.
---

# Output Assert DSL

Package: `github.com/xhd2015/doctest/assert`

Readable **template-shaped expected output** for doctest `ASSERT.md` leaves. Replaces ad-hoc `strings.Contains`, `strings.Index`, and hand-rolled parsing.

Design history: `doc/DESIGN_OUTPUT_ASSERT.md`. Older dialects:

| Version | Package / status |
|---------|------------------|
| **v3 (default)** | root `assert` — preferred for new tests |
| **v2** | `assert/legacy_v2` — **deprecated**; only with explicit `version: 2` |
| **v1** | `assert/legacy_v1` — tag DSL (`<contains>`, …) when not using YAML dialect |

## When to use

| Situation | Approach |
|-----------|----------|
| Bounded stdout/stderr (typical leaf) | `assert.Output(t, actual, template)` — strict full match |
| Variable port, path, user, count | `__PLACEHOLDER__` in YAML header + template body |
| Middle section may differ (stack trace, progress) | `...N lines omitted...` |
| Flexible single line (dots, prefixes) | raw regex on the line (e.g. `^\.+$`) |
| Literal dots / parens in output | Escape: `0\.001s`, `foo\(bar\)` |
| ANSI-colored segments | `<ansi-color bold gray>…</ansi-color>` (inner text is literal) |
| Structured variable (SHA, UUID) | `regex=` on the placeholder def |
| Document example values for readers | `example=…` in placeholder header (metadata only) |

**v3 is strict line-by-line full match.** There is no `<contains>` or `assert.Contains()` in v3 templates.

## API

```go
import "github.com/xhd2015/doctest/assert"

p, err := assert.Parse(template)
err = assert.Match(p, actual)       // strict full match (default)
assert.Output(t, actual, template)  // Parse + Match + t.Fatal
```

### Version routing

| Template | Engine |
|----------|--------|
| YAML header with `version: 3` | **v3** |
| YAML header with placeholders / YAML dialect and **no** version key | **v3** (default) |
| YAML header with `version: 2` | **legacy_v2** (deprecated) |
| Unknown `version:` value | **parse error** |
| v1 tag DSL (no YAML version dialect) | **legacy_v1** |

- `\r\n` → `\n` normalization is always on.
- Trailing newlines are **strict** (template and actual must agree).
- Matching is **case-sensitive**.

### CLI stdout trailing newline (required)

User-facing CLI stdout **must** end with a final `\n` (POSIX convention — keeps
the shell prompt on its own line). v3 templates for CLI output **must** include
that trailing newline too.

**Go authoring pattern** — put the closing backtick on the line **after** the
last content line (not on the same line, and not separated by an extra blank
line):

````go
assert.Output(t, resp.Stdout, `---
version: 3
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

## Template shape (v3)

```DSL
---
version: 3
__PORT__: type=number, example=8901, a port
__USER__: type=string, example=alice, logged-in user
__SHA__: regex=[0-9a-f]{7,40}, example=abc1234, short git sha
---
Server listen on: __PORT__
...3 lines omitted...
Hello __USER__
commit __SHA__
<ansi-color bold gray>1 Cached</ansi-color>
```

| Part | Role |
|------|------|
| `---` … `---` | YAML header; omit version or set `version: 3` for v3; `version: 2` is deprecated |
| `__NAME__:` | Placeholder definition (`type=` and/or `regex=`, plus `example=` / explanation) |
| Body lines | Each **content** line is a **raw Go regular expression**, full-line match (`^…$`) |

In Go `ASSERT.md` blocks, pass the template as one raw string starting with `---` — no `` + `` concatenation needed. Leading blank lines before the opening `---` are trimmed; trailing blank lines in the body are preserved.

### Placeholder definitions

Compact form (preferred):

```yaml
__PORT__: type=number, example=8901, a port
__SHA__: regex=[0-9a-f]{7,40}, example=a1b2c3d, short sha
```

- `k=v` pairs before the first bare word — machine metadata
- Trailing text after metadata — human explanation only (ignored by matcher)

| Field | Required | Values |
|-------|----------|--------|
| `type` | one of `type` or `regex` | `string`, `number` |
| `regex` | one of `type` or `regex` | Go regexp **fragment** (subpattern) |
| `example` | no | documentation for readers |

| Spec | Matches on one line |
|------|---------------------|
| `type=string` | any non-newline text (`[^\n]*?`, non-greedy) |
| `type=number` | integer or float (`-?\d+(?:\.\d+)?`) — **loose** |
| `regex=…` | custom subpattern; invalid fragment → **parse error** |

**Do not set both `type=` and `regex=`** on the same placeholder — that is a **parse error**.

Every `__NAME__` used in the body must be defined in the header.

Placeholders expand to **named capture groups** (`(?P<NAME>…)`). If the same
`__NAME__` appears more than once, all captures must be the **same string**
(same-value binding). Mismatches fail the match with a clear error.

### Body line kinds

1. **Content line (default)** — entire line is a Go regexp matched as a full line.
   - Metacharacters are active: `.` matches any character, `(` starts a group, etc.
   - Escape literals with `\`: write `0\.001s` for the text `0.001s`.
   - `__PLACEHOLDER__` is replaced by the placeholder subpattern (named group).
   - There is **no** separate “regex intent” scan and **no** automatic QuoteMeta of the whole line.

2. **Omit marker** — `...N lines omitted...` (whitespace allowed around parts).
   Special (not a content regex). Skips exactly **N** actual lines (any content).
   `N` must be a non-negative integer.

### Escaping cheatsheet (content lines)

| Want literal | Write |
|--------------|--------|
| `0.001s` | `0\.001s` |
| `foo(bar)` | `foo\(bar\)` |
| `a+b` | `a\+b` |
| `src/main.go` | `src/main\.go` (or a `__PATH__` placeholder) |
| Progress dots only | `^\.+$` (intentional regex) |

### Color spans

`<ansi-color SPEC>inner text</ansi-color>` — only inline structural tag in v3.

- Open SGR + reset envelope as today.
- **Inner text is QuoteMeta’d** (literal UI text). You do **not** need `\.` inside the tag for a literal dot.
- Outside the tag, the rest of the line remains raw regex.

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
version: 3
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

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatal(err)
    }
    assert.Output(t, resp.Stdout, `---
version: 3
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

### Regex is per line only

Content regex never crosses newlines. Use `...N lines omitted...` to skip variable middle sections instead of multiline regex.

### Same-value binding

```DSL
---
__ID__: type=string, example=abc
---
start __ID__
end __ID__
```

Both `__ID__` captures must be identical. Different values → match error.

## Migration notes

### From ad-hoc Go

| Legacy pattern | v3 prefer |
|----------------|-----------|
| `strings.Contains` loop | strict template with `...N lines omitted...` or a flexible regex line |
| `strings.Index` + `strings.Count` for dots | `^\.+$` |
| Dual `Contains` for platform errors | `(linux-msg\|darwin-msg)` alternation |
| `metricIsColored` / `stripANSI` helpers | `<ansi-color bold gray>…</ansi-color>` |

### From v2

| v2 habit | v3 |
|----------|-----|
| Literal dots without escaping (`0.001s`) | Escape: `0\.001s` |
| Dual “pattern vs regex-intent” lines | Every content line is raw regex |
| `version: 2` required for YAML dialect | Default is v3; omit version or set `version: 3` |
| No capture equality | Repeated `__NAME__` must match the same value |
| Only `type=string` / `type=number` | Plus `regex=` custom fragments |

Existing v1 templates continue to work (no YAML version dialect). Prefer **v3** for new tests. **v2 is deprecated** (`assert/legacy_v2`); do not start new leaves on `version: 2`.

## Real-world CLI cookbook

CLI templates in this section should end with a trailing `\n` in both the
template and the simulated actual bytes. When authoring `ASSERT.md` Go blocks,
use the closing-backtick-on-next-line pattern from **CLI stdout trailing newline**
above.

**188** doctest leaves under `assert/tests/output-assert-v3-suite/integration/real-world/`
(17 categories). All use **simulated** bytes — no subprocess. Regenerate via
`go run ./script/generate/real-world-assert-cases/main.go`.

Categories: `unix-text`, `go-toolchain`, `rust-toolchain`, `node-js`, `python`,
`git`, `http-clients`, `containers`, `build-systems`, `databases`, `jvm-kotlin`,
`c-cpp`, `shell`, `package-managers`, `cloud-infra`, `languages-other`,
`misc-devtools`.

**YAML tip:** quote placeholder defs when `example=` contains colons:
`__LINE__: 'type=string, example=go: creating module foo'`.

Samples below; every leaf in the tree is a copy-pasteable template.

### cat — file dump (escape dots)

```DSL
---
version: 3
---
# My Project
Version 1\.0
```

### grep -n — line number + match

```DSL
---
version: 3
__LINE__: type=number, example=3, 1-based line number
---
__LINE__:func main\(\) \{
```

### rg — path:line:column (no heading)

Use one placeholder for the variable prefix when literals follow on the same line
(string placeholders are line-greedy).

```DSL
---
version: 3
__HIT__: type=string, example=src/main.go:42:5, ripgrep hit prefix
---
__HIT__
```

### go build — compile error with omitted notes

```DSL
---
version: 3
__PACKAGE__: type=string, example=example.com/foo, module import path
---
# __PACKAGE__
\./main\.go:10:2: undefined: Bar
...2 lines omitted...
FAIL
```

### go test — pass summary with colored PASS

Keep stable substrings escaped when needed; parameterize only the trailing timing field.

```DSL
---
version: 3
__SECONDS__: type=number, example=0.123, elapsed seconds
---
=== RUN   TestFoo
--- PASS: TestFoo \(0\.00s\)
<ansi-color green>PASS</ansi-color>
ok  	example\.com/foo	__SECONDS__s
```

### go mod init

```DSL
---
version: 3
__MODULE__: type=string, example=example.com/myproject, new module path
---
go: creating new go\.mod: module __MODULE__
```

### npm run build — script banner + omitted log + timing

Combine name@version into one placeholder when followed by more text on the line.

```DSL
---
version: 3
__BANNER__: type=string, example=myapp@1.0.0 build, lifecycle banner
__SECONDS__: type=number, example=1.2, build duration
---
> __BANNER__
> tsc
...4 lines omitted...
Done in __SECONDS__s\.
```

### npm init — path + omitted JSON body

```DSL
---
version: 3
__PATH__: type=string, example=/tmp/proj/package.json, output path
---
Wrote to __PATH__
...4 lines omitted...
\}
```

### curl -i — HTTP headers + omitted body

```DSL
---
version: 3
__CODE__: type=number, example=200, HTTP status code
---
HTTP/1\.1 __CODE__ OK
Content-Type: application/json
...3 lines omitted...
\{"status":"ok"\}
```

## Examples in repo

```sh
doctest test ./assert/tests/output-assert-v3/...          # focused v3 engine suite
doctest test ./assert/tests/output-assert-v3-suite/...    # full suite + 188 CLI leaves
doctest test ./assert/tests/output-assert/...             # v1 legacy
go test ./assert/...
```
