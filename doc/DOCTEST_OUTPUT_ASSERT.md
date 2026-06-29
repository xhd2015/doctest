---
name: doctest-output-assert
description: >-
  Output assert DSL for doctest ASSERT.md — template tags, API, and migration
  from strings.Contains. Use when asserting CLI or text output in doc-style tests.
---

# Output Assert DSL

Package: `github.com/xhd2015/doctest/assert`

Readable **template-shaped expected output** for doctest `ASSERT.md` leaves. Replaces ad-hoc `strings.Contains`, `strings.Index`, and hand-rolled parsing.

Full design history: `doc/DESIGN_OUTPUT_ASSERT.md` in the doctest repo.

## When to use

| Situation | Approach |
|-----------|----------|
| Bounded stdout/stderr (typical leaf) | `assert.Output(t, actual, template)` — full exact match |
| Bounded output with scattered required lines | ``assert.Output(t, actual, `<contains>...</contains>`)`` |
| Excerpt inside a long log | `assert.Match(p, actual, assert.Contains())` — contiguous subregion |
| Help keywords, scattered lines | `<contains>` + `<start-with>` / `<end-with>` |
| Platform-specific error text | `<any-of><expect>…</expect></any-of>` |
| Variable progress dots | `<regex>^\.+$</regex>` |
| ANSI-colored segments | `<ansi-color bold gray>…</ansi-color>` |
| Annotate literal path/id for readers | `<hint:path>~/proj</hint:path>` (still literal match) |

## API

```go
import "github.com/xhd2015/doctest/assert"

p, err := assert.Parse(template)
err = assert.Match(p, actual)                      // MatchExact (default)
err = assert.Match(p, actual, assert.Contains())   // contiguous subregion
assert.Output(t, actual, template)                 // Parse + Match + t.Fatal
```

- `\r\n` → `\n` normalization is always on.
- Trailing newlines are **strict** (template and actual must agree).
- Matching is **case-sensitive**.

## ASSERT.md pattern

**Recommended** (not required by `doctest vet`): mirror the template in prose:

````markdown
## Expected Output

```
<contains>
Usage: mytool
<start-with>
  build
</start-with>
</contains>
```

## Expected
- stdout lists build subcommand

```go
import (
    "testing"
    "github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatal(err)
    }
    assert.Output(t, resp.Stdout, `
<contains>
Usage: mytool
<start-with>
  build
</start-with>
</contains>`)
}
```
````

## Included tags

Authoritative registry. **No other tag forms are supported.**

### Summary table

| Tag | Form | Consumes actual output? | Matching |
|-----|------|-------------------------|----------|
| `<optional>` | block **or** inline | block meta: **no**; inline inner: optional | Block: inner lines absent or present; inline: wrapped span absent or present |
| `<any-of>` | block **or** inline | block meta: **no**; inline: **yes** (chosen branch) | Exactly one `<expect>` branch must match |
| `<expect>` | block or inline (inside `<any-of>`) | block meta: **no** | Delimits one branch; body matched literally |
| `<regex>` | block **or** inline | **yes** — matched line/span | Go regexp full match |
| `<contains>` | block only | **no** (meta) | Every inner fragment found in actual (any order); default **full line** |
| `<start-with>` | inline or block (inside `<contains>`) | **no** (modifier) | Fragment must match line prefix |
| `<end-with>` | inline or block (inside `<contains>`) | **no** (modifier) | Fragment must match line suffix |
| `<hint:label>` | inline only | **yes** — inner text | Inner text must match literally; label is documentation |
| `<ansi-color>` | inline only | **yes** — inner text | Inner text literal; strict ANSI wrapper |

### Block meta tags (standalone lines)

Tags alone on their line (whitespace only besides). Never appear in CLI output; never consume an actual output line (except `<regex>` inner pattern line).

| Open | Close | Purpose |
|------|-------|---------|
| `<optional>` | `</optional>` | Optional multiline region |
| `<any-of>` | `</any-of>` | Alternative branches |
| `<expect>` | `</expect>` | One branch inside `<any-of>` |
| `<contains>` | `</contains>` | Order-free fragments |
| `<regex>` | `</regex>` | One-line pattern region (block form) |

### Inline tags (same line as content)

| Tag | Example | Purpose |
|-----|---------|---------|
| `<optional>…</optional>` | `Result: <optional>warn: </optional>OK` | Only wrapped span is optional |
| `<hint:label>…</hint:label>` | `id=<hint:id>abc-123</hint:id>` | Annotate literal span |
| `<ansi-color SPEC>…</ansi-color>` | `<ansi-color bold gray>1 Cached</ansi-color>` | Assert ANSI style (strict) |
| `<any-of>…</any-of>` | `status: <any-of><expect>ok</expect><expect>done</expect></any-of>` | Inline alternatives |
| `<expect>…</expect>` | inside inline `<any-of>` | One inline branch |
| `<regex>…</regex>` | `<regex>^\.+$</regex>` | Regexp on line or span |
| `<start-with>…</start-with>` | inside `<contains>` | Line prefix match |
| `<end-with>…</end-with>` | inside `<contains>` | Line suffix match |

Ordinary lines inside `<contains>` may also use inline pattern tags such as
`<any-of>`, `<optional>`, `<regex>`, and `<hint:...>`.

### Reserved names and invalid forms

| Valid | Invalid (rejected at parse) |
|-------|----------------------------|
| `<hint:id>…</hint:id>` | `<id>…</id>` — bare name, no `hint:` prefix |
| `<hint:path>…</hint:path>` | `<cached>…</cached>` — not a registered tag |
| `<optional>…</optional>` | `<>` / `</>` — removed syntax |
| `<ansi-color #90>…</ansi-color>` | `<ansi-color orange>…` — unknown name without `#` |

**Registered top-level names:** `optional`, `any-of`, `expect`, `contains`, `regex`, `start-with`, `end-with`, `ansi-color`, and `hint` (always with `:label` suffix).

### `<ansi-color>` style specifier

Space-separated tokens; **strict** open SGR + `\x1b[0m` reset immediately after inner text.

| Token | Open SGR | Example |
|-------|----------|---------|
| `bold` | `\x1b[1m` | `<ansi-color bold>title</ansi-color>` |
| `red` | `\x1b[31m` | `<ansi-color red>FAIL</ansi-color>` |
| `green` | `\x1b[32m` | `<ansi-color green>OK</ansi-color>` |
| `gray` | `\x1b[90m` | `<ansi-color gray>1 Cached</ansi-color>` |
| `#` + params | `\x1b[<params>m` | `<ansi-color #38;5;208>warn</ansi-color>` |

Combined — emitted left to right: `<ansi-color bold gray>1 Cached</ansi-color>` → `\x1b[1m\x1b[90m` + text + reset.

### What is NOT a tag

| Construct | Status |
|-----------|--------|
| `<run>`, `<pass>`, `<cached>`, `<cwd>` | **Rejected** |
| `<dots>…</dots>` | **Rejected** — use `<regex>^\.+$</regex>` |
| User-defined `<foo>…</foo>` | **Rejected** |
| Free regex in literals | **Rejected** — use `<regex>` tag |

### Quick reference

```
<optional>              ← block meta
  optional lines
</optional>

<any-of>                ← block meta
<expect>
branch A
</expect>
<expect>
branch B
</expect>
</any-of>

<regex>
^\.+$
</regex>

<contains>
Usage: mytool
<start-with>
  build
</start-with>
</contains>

  (<ansi-color gray>1 Cached</ansi-color>, 0 Fail)
Result: <optional>warn: </optional>OK
$ cd <hint:path>~/proj</hint:path>
```

## Key semantics

### Optional (none or full match)

Block `<optional>` is **absent** (zero lines) or inner body **fully** matches. Adjacent block optionals are **separate, never merged**.

### Any-of failure

On failure, report **all branches** with want/got per branch.

### Contains vs `assert.Contains()` option

| Mechanism | Order | Adjacency |
|-----------|-------|-----------|
| `<contains>` tag | Order-free between fragments | Fragments anywhere in actual |
| `assert.Contains()` option | Order preserved | Entire template as one contiguous slice |

Prefer ``assert.Output(t, actual, `<contains>...</contains>`)`` when the template
itself is a `<contains>` block. Avoid combining a top-level `<contains>` block
with `assert.Match(..., assert.Contains())`: that mixes order-free fragment
matching with contiguous excerpt matching and usually indicates the assertion is
over-specified or unclear.

### Escaping in templates

Most `<` `>` in output need no escape. Use `\<tag>` / `\</tag>` only when text would parse as a registered tag. Backslash is template-only (not in actual).

### Hints

`<hint:label>…</hint:label>` — inner text matches **literally**; label is for readers and errors only.

## Migration from legacy asserts

| Legacy pattern | Prefer |
|----------------|--------|
| `strings.Contains` loop over stdout lines | `<contains>` + `<start-with>` |
| `strings.Index` + `strings.Count` for dots | `<regex>^\.+$</regex>` |
| Dual `Contains` for platform errors | `<any-of><expect>…</expect></any-of>` |
| `metricIsColored` / `stripANSI` helpers | `<ansi-color bold gray>…</ansi-color>` |

Existing trees may keep legacy asserts; doctest review flags them as **suggestions** to migrate.

## Examples in repo

```sh
go test ./assert/...
doctest test -v ./assert/tests/...
```
