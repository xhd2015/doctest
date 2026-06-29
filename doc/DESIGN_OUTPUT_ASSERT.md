# Output Assert DSL — Design Document

> **Status:** **Implemented** (`github.com/xhd2015/doctest/assert`). This document is the canonical DSL reference.
>
> **Package:** `github.com/xhd2015/doctest/assert`
>
> **Goal:** A readable template language for asserting CLI (and similar text) output in doctest `ASSERT.md` leaves, replacing ad-hoc `strings.Contains` / hand-rolled parsing.
>
> **Skill:** `doctest skill output-assert show`
>
> ## Authoring quick start
>
> ```go
> import "github.com/xhd2015/doctest/assert"
>
> func Assert(t *testing.T, req *Request, resp *Response, err error) {
>     // ...
>     assert.Output(t, resp.Stdout, `
> <contains>
> Usage: mytool
> <start-with>
>   build
> </start-with>
> </contains>`)
> }
> ```
>
> - Default: **`assert.Output`** = full-output exact match.
> - Long logs: **`assert.Match(p, actual, assert.Contains())`** for a contiguous excerpt.
> - Tag reference: **§3** below.
> - Optional prose mirror: `## Expected Output` fenced block in ASSERT.md (advertised, not required by vet — §16).

---

## 1. Problem

Doctest leaves assert CLI output today with imperative Go:

```go
if !strings.Contains(resp.Summary, "1 Cached") { ... }
summaryIdx := strings.Index(resp.Output, "  (")
dotsBefore := strings.Count(resp.Output[:summaryIdx], ".")
```

This works but scales poorly:

| Approach | Strength | Weakness |
|----------|----------|----------|
| `strings.Contains` | Easy to write | Too loose; passes when unrelated text appears |
| Exact string equality | Strict | Brittle on timestamps, paths, PIDs, cache counts |
| Hand-rolled parsing | Precise | Verbose, hard to read, duplicated per test |
| Raw `regexp` | Powerful | Hard to maintain; obscures intent |

We want something closer to **PHP-embedded HTML** or **annotated golden output**: mostly literal text the human can read, with small, explicit escape hatches for non-deterministic regions.

---

## 2. Core Idea

A **template-shaped expected output** with five constructs (see **§3 Included Tags** for the full registry):

| Kind | Syntax | Meaning |
|------|--------|---------|
| **Literal** | Everything outside tags | Byte-exact match (modulo newline/ANSI policy — see §6) |
| **Optional** | `<optional>…</optional>` | Region that may be absent; see §2.2 for block vs inline |
| **Any-of** | `<any-of><expect>…</expect>…</any-of>` | Match **one** of several alternative branches (block or inline) |
| **Regex** | `<regex>…</regex>` | Line or span matches Go regexp (block or inline) |
| **Contains** | `<contains>…</contains>` | Every inner line must appear in actual output (order-free) |
| **Hint** | `<hint:label>…</hint:label>` | Documentation marker only; **inner text still matches literally** |

### 2.1 Literal text

Default mode. What you see is what must appear.

```
Usage: doctest [command] [options]
```

> Only **tag-shaped** `<…>` / `</…>` needs escaping in templates (`\<tag>` — §5.10). Plain `2 < 3` or `[command]` need no escape.

### 2.2 Optional — block vs inline

Optional content uses `<optional>…</optional>`. The same tag name serves two forms with **different matching scope**:

#### Block form (standalone meta lines)

Opening/closing tags occupy **their own lines** and are **template scaffolding** — they never appear in CLI output and **do not consume** actual output lines (including the newline after `<optional>`).

```
Usage: doctest [command] [options]

Commands:
  build    Build test binaries
  test     Run tests
<optional>
  agent    Agent commands (experimental)
</optional>
  skill    Skill management
```

| Rule | Detail |
|------|--------|
| Absent | Inner lines not present in actual → still matches |
| Present | Inner lines match literally (or contain inline hints/optionals) |
| Position | Must sit between surrounding literal anchors; order matters |
| Meta lines | `<optional>` and `</optional>` lines are not matched against output |

#### Inline form (same line as content)

When `<optional>` and `</optional>` wrap text **on the same line**, only the wrapped span is optional; the rest of the line is literal.

```
Result: <optional>warning: </optional>OK
```

| Actual | Match? |
|--------|--------|
| `Result: OK` | yes — optional span absent |
| `Result: warning: OK` | yes — optional span present |
| `OK` | no — `Result: ` prefix is literal |

#### Block vs inline disambiguation

| Form | Detection |
|------|-----------|
| **Block** | `<optional>` is the only non-whitespace on its line (same for `</optional>`) |
| **Inline** | `<optional>` shares a line with literal or hint text before `</optional>` |

### 2.3 Any-of — alternatives

For output that may legitimately differ (platform messages, pass vs fail summaries):

```
<any-of>
<expect>
error: file not found
</expect>
<expect>
error: permission denied
</expect>
</any-of>
```

- `<any-of>`, `</any-of>`, `<expect>`, `</expect>` on **standalone lines** are **meta tags** — same rule as block `<optional>`: they do not consume actual output lines.
- Matcher succeeds if **any one** `<expect>…</expect>` branch matches the next region of actual output.
- Branches are tried in source order; first match wins on success; on failure **all branches** are reported (§12 Q6).

**Inline any-of** (§12 Q7): when `<any-of>` shares a line with literals, only the wrapped alternative span is variable; prefix/suffix on the same line remain literal — same scoping rule as inline `<optional>`.

```
status: <any-of><expect>ok</expect><expect>done</expect></any-of>
```

| Actual | Match? |
|--------|--------|
| `status: ok` | yes |
| `status: done` | yes |
| `ok` | no — `status: ` is literal |

### 2.4 Hints — documentation markers (literal match)

Hints annotate regions of output for human readers and clearer error messages. They are **not** wildcards — the text inside must **literally match** actual output.

```
$ cd <hint:path>~/Projects/myapp</hint:path>
  (2 Run, 2 Pass, 1 Cached, 0 Fail)
  Took <hint:duration>1.2s</hint:duration>
```

| Template | Actual | Result |
|----------|--------|--------|
| `<hint:path>~/Projects/myapp</hint:path>` | `~/Projects/myapp` | pass |
| `<hint:path>~/Projects/myapp</hint:path>` | `/tmp/other` | fail — hint does not relax matching |
| `<hint:id>abc-123</hint:id>` | `abc-123` | pass |

- **Syntax:** `<hint:label>…</hint:label>` — `label` is a short semantic name (`id`, `path`, `duration`, …).
- **Prefix required:** bare `<id>`, `<path>`, `<cached>` are **not valid**; always use the `hint:` prefix.
- **Purpose:** readability + failures like `hint:path expected "~/Projects/myapp", got "/tmp/other"`.
- **Non-goals:** hints do not accept alternate values. Use `<any-of>`, `<optional>`, or `<contains>` for variation.

**Omit hints when not useful.** Summary fields like `1 Cached` need no annotation — write the literal:

```
  (2 Run, 2 Pass, 1 Cached, 0 Fail)
```

When counts legitimately differ, use `<any-of>` with separate literal branches (§4.1), not a hint.

### 2.5 Contains — order-free fragments

Block `<contains>` replaces `strings.Contains` loops. Each inner fragment must appear **somewhere** in actual output; order does not matter. Meta lines do not consume output (§3).
Ordinary fragment lines may use inline pattern tags such as `<any-of>`,
`<optional>`, `<regex>`, and `<hint:...>`.

- **Default:** fragment is a **full line** match.
- **`<start-with>` / `<end-with>`:** prefix or suffix match when a full line is too strict.

```
<contains>
Usage: doctest
<start-with>
  agent
</start-with>
</contains>
```

When the template itself is a top-level `<contains>` block, prefer:

```go
assert.Output(t, actual, `<contains>
Usage: mytool
</contains>`)
```

Do not combine a top-level `<contains>` block with
`assert.Match(p, actual, assert.Contains())`; that mixes order-free fragment
matching with contiguous excerpt matching. Use `assert.Contains()` only when the
template is a contiguous excerpt from a larger output.

---

## 3. Included Tags

Authoritative list of tags in the DSL. **No other tag forms are supported in MVP.**

### 3.1 Summary table

| Tag | Form | Consumes actual output? | Matching |
|-----|------|-------------------------|----------|
| `<optional>` | block **or** inline | block meta: **no**; inline inner: optional | Block: inner lines absent or present; inline: wrapped span absent or present |
| `<any-of>` | block **or** inline | block meta: **no**; inline: **yes** (chosen branch) | Exactly one `<expect>` branch must match |
| `<expect>` | block (inside block `<any-of>`) or inline (inside inline `<any-of>`) | block meta: **no** | Delimits one branch; body matched literally |
| `<regex>` | block **or** inline | **yes** — matched line/span | Go regexp full match (§5.9) |
| `<contains>` | block only | **no** (meta) | Every inner fragment found in actual (any order); default **full line** |
| `<start-with>` | inline or block (inside `<contains>`) | **no** (modifier) | Fragment must match line prefix |
| `<end-with>` | inline or block (inside `<contains>`) | **no** (modifier) | Fragment must match line suffix |
| `<hint:label>` | inline only | **yes** — inner text | Inner text must match literally; label is documentation |
| `<ansi-color>` | inline only | **yes** — inner text | Inner text literal; must be wrapped in named or raw ANSI color |

### 3.2 Block meta tags (standalone lines)

These tags appear **alone on their line** (whitespace only besides). They never appear in CLI output and never consume an actual output line.

| Open | Close | Purpose |
|------|-------|---------|
| `<optional>` | `</optional>` | Optional multiline region |
| `<any-of>` | `</any-of>` | Alternative branches |
| `<expect>` | `</expect>` | One branch inside `<any-of>` |
| `<contains>` | `</contains>` | Order-free substring lines |
| `<regex>` | `</regex>` | One-line pattern region (block form) |

### 3.3 Inline tags (same line as content)

| Tag | Example | Purpose |
|-----|---------|---------|
| `<optional>…</optional>` | `Result: <optional>warn: </optional>OK` | Only wrapped span is optional |
| `<hint:label>…</hint:label>` | `id=<hint:id>abc-123</hint:id>` | Annotate literal span |
| `<ansi-color SPEC>…</ansi-color>` | `<ansi-color bold gray>1 Cached</ansi-color>` | Assert span has ANSI style (strict) |
| `<any-of>…</any-of>` | `status: <any-of><expect>ok</expect><expect>done</expect></any-of>` | Inline alternatives |
| `<expect>…</expect>` | inside inline `<any-of>` | One inline branch |
| `<regex>…</regex>` | `<regex>^\.+$</regex>` | Regexp match on line or span |
| `<start-with>…</start-with>` | inside `<contains>` | Line prefix match |
| `<end-with>…</end-with>` | inside `<contains>` | Line suffix match |

### 3.4 Reserved names and invalid forms

| Valid | Invalid (rejected at parse) |
|-------|----------------------------|
| `<hint:id>…</hint:id>` | `<id>…</id>` — bare name, no `hint:` prefix |
| `<hint:path>…</hint:path>` | `<cached>…</cached>` — not a registered tag |
| `<optional>…</optional>` | `<>` / `</>` — removed syntax |
| `<ansi-color #90>…</ansi-color>` | `<ansi-color orange>…` — unknown name without `#` |

**Registered top-level names:** `optional`, `any-of`, `expect`, `contains`, `regex`, `start-with`, `end-with`, `ansi-color`, and `hint` (always with `:label` suffix).

**`<ansi-color>` style specifier** — space-separated tokens, **strict** open SGR + reset (§5.12, §12 Q15):

| Token | Open SGR | Example |
|-------|----------|---------|
| `bold` | `\x1b[1m` | `<ansi-color bold>title</ansi-color>` |
| `red` | `\x1b[31m` | `<ansi-color red>FAIL</ansi-color>` |
| `green` | `\x1b[32m` | `<ansi-color green>OK</ansi-color>` |
| `gray` | `\x1b[90m` | `<ansi-color gray>1 Cached</ansi-color>` |
| `#` + params | `\x1b[<params>m` | `<ansi-color #38;5;208>warn</ansi-color>` |

Combined styles — tokens emitted **left to right**:

```
<ansi-color bold gray>1 Cached</ansi-color>
```

→ `\x1b[1m\x1b[90m` + `1 Cached` + `\x1b[0m` (strict adjacency).

**Disambiguation:** `#…` is one token. Other tokens must be registered names (`bold`, `red`, `green`, `gray`). Unknown tokens are parse errors.

**`<regex>` pattern:** Go `regexp` (RE2) syntax. Pattern text is literal in the template (not a wildcard tag).

**`<ansi-color>` / `<regex>` close tags:** always `</ansi-color>` / `</regex>` (no specifier on close).

**Hint labels:** `[a-zA-Z][a-zA-Z0-9_]*` — e.g. `id`, `path`, `duration`. Open/close must use the same label.

### 3.5 What is NOT a tag

| Construct | Status |
|-----------|--------|
| `<run>`, `<pass>`, `<cached>`, `<cwd>` | **Rejected** — use literals, `<any-of>`, or `<hint:label>` |
| `<dots>…</dots>` | **Rejected** — use `<regex>^\.+$</regex>` |
| User-defined `<foo>…</foo>` | **Rejected** — only §3.1 tags allowed |
| Free regex in literals | **Rejected** — use `<regex>` tag only |
| Typed tags `<hint:id:uuid>` | **Future** — not MVP |

### 3.6 Quick reference

```
<optional>              ← block meta
  optional lines
</optional>

<any-of>                ← block meta
<expect>               ← block meta
branch A
</expect>
<expect>
branch B
</expect>
</any-of>

<regex>                ← block meta (pattern line consumes one actual line)
^\.+$
</regex>

status: <any-of><expect>ok</expect><expect>done</expect></any-of>  ← inline any-of

<contains>             ← block meta
Usage: doctest          ← full-line fragment (default)
<start-with>
  agent
</start-with>          ← prefix match inside contains
build
</contains>

  (<ansi-color gray>1 Cached</ansi-color>, <ansi-color #90>0 Fail</ansi-color>)  ← named or raw ANSI

Result: <optional>warn: </optional>OK          ← inline optional
$ cd <hint:path>~/proj</hint:path>               ← inline hint (literal match)
  (2 Run, 2 Pass, 1 Cached, 0 Fail)             ← plain literals (no tag needed)
```

---

## 4. Motivating Examples

### 4.1 Doctest build summary (replaces color/cached asserts)

**Today** (`libdoc/build/tests/output-color/.../ASSERT.md`):

```go
if !strings.Contains(resp.Summary, "1 Cached") { ... }
if !metricIsColored(resp.Summary, "1 Cached") { ... }
```

**With DSL** — literals for fixed counts; `<any-of>` when cached count varies:

```
..
  (2 Run, 2 Pass, 1 Cached, 0 Fail)
```

Or when either cached or uncached summary is acceptable:

```
<any-of>
<expect>
  (2 Run, 2 Pass, 1 Cached, 0 Fail)
</expect>
<expect>
  (2 Run, 2 Pass, 0 Cached, 0 Fail)
</expect>
</any-of>
```

No `<hint:cached>` — not useful when the literal `1 Cached` / `0 Cached` already states the expectation.

**ANSI color** (replaces `metricIsColored` helper):

```
..
  (2 Run, 2 Pass, <ansi-color gray>1 Cached</ansi-color>, <ansi-color gray>0 Fail</ansi-color>)
```

Asserts `1 Cached` and `0 Fail` appear literally and are wrapped in gray ANSI codes — no separate Go helper.

### 4.2 Dot progress (replaces manual dot counting)

**Today** (`libdoc/build/tests/dot-progress/incremental/ASSERT.md`): ~40 lines of `strings.Index`, `strings.Count`, loop over trailing lines.

**With DSL** — `<regex>` for variable dot count (§12 Q11); literal summary line:

```
<regex>
^\.+$
</regex>
  (2 Run, 2 Pass, 0 Cached, 0 Fail)
```

Or inline on one template line:

```
<regex>^\.+$</regex>
  (2 Run, 2 Pass, 0 Cached, 0 Fail)
```

| Tag | Role | Why |
|-----|------|-----|
| `<regex>^\.+$</regex>` | One actual line of one-or-more `.` chars | Package count varies → dot count varies |
| `(2 Run, 2 Pass, 0 Cached, 0 Fail)` | Literals | Counts fixed by test setup |

Optional wrapper if dots may be absent entirely:

```
<optional>
<regex>
^\.+$
</regex>
</optional>
  (2 Run, 2 Pass, 0 Cached, 0 Fail)
```

### 4.3 Help text (replaces `strings.Contains` loops)

**Today** (`tests/help/top-level/ASSERT.md`):

```go
for _, want := range []string{"Usage: doctest", "agent", "build", ...} {
    if !strings.Contains(resp.Stdout, want) { ... }
}
```

**With DSL** — two viable styles:

**Exact block** (strict ordering):

```
Usage: doctest [command] [options]

Commands:
  agent    Agent commands
  build    Build test binaries
  test     Run tests
  skill    Skill management
```

**Contains block** (order-free — see §3.1 `<contains>`). Command names use `<start-with>` because help lines are indented:

```
<contains>
Usage: doctest
<start-with>
  agent
</start-with>
<start-with>
  build
</start-with>
<start-with>
  test
</start-with>
<start-with>
  skill
</start-with>
</contains>
```

### 4.4 Platform-specific error text (any-of)

**With DSL:**

```
<any-of>
<expect>
Error: file already closed
</expect>
<expect>
Error: bad file descriptor
</expect>
</any-of>
```

Replaces dual `strings.Contains` checks in `libdoc/cli/tests/read-stdin-error/.../ASSERT.md`.

---

## 5. Syntax Specification (Draft)

### 5.1 Line orientation

CLI output is line-based. The matcher operates primarily on lines (`\n`-terminated rows).

- A template is a sequence of **lines**, **block meta regions**, and **inline segments**.
- A single line may contain **inline** literal / hint / optional spans.
- Block meta regions (`<optional>`, `<any-of>`, `<contains>`) span multiple template lines but contribute **zero** lines to output matching on their own.

### 5.2 Block meta tags vs inline tags

Structural tags split into two roles:

| Role | Tags | Line shape | Consumes actual output? |
|------|------|------------|-------------------------|
| **Block meta** | `<optional>`, `<any-of>`, `<expect>`, `<contains>`, `<regex>` (+ closers) | Tag alone on its line (only whitespace besides) | **No** for scaffolding; **yes** for inner pattern line of `<regex>` |
| **Inline structural** | `<optional>…</optional>`, `<any-of>…</any-of>` | Tag shares a line with matched text | Only wrapped span is optional / any-of |
| **Hint** | `<hint:label>…</hint:label>` | Always inline | Inner text matched **literally**; label for docs/errors |
| **Regex** | `<regex>…</regex>` | Block or inline | Pattern full-matches one line (block/line) or span (inline) |

**Key rule (from design):** A standalone meta-tag line (e.g. `<any-of>`) and the newline after it **do not consume** any line from the target output. Only the *inner* template lines between meta open/close are candidates for matching.

```
Template                          Actual
────────                          ──────
<any-of>                    →     (nothing consumed)
<expect>                    →     (nothing consumed)
error: file not found       →     matches this line
</expect>                   →     (nothing consumed)
</any-of>                   →     (nothing consumed)
OK                          →     matches this line
```

### 5.3 Grammar (informal)

```
Template       := Item*
Item           := LiteralLine | PatternLine | BlockOptional | AnyOfBlock | ContainsBlock | BlockRegex

LiteralLine    := text without unescaped tags, newline
PatternLine    := (Literal | InlineOptional | InlineAnyOf | Hint | AnsiColor | InlineRegex)+ newline

Hint           := '<hint:' Label '>' Text '</hint:' Label '>'
AnsiColor      := '<ansi-color' StyleSpec '>' Text '</ansi-color>'
StyleSpec      := StyleToken+
StyleToken     := 'bold' | 'red' | 'green' | 'gray' | '#' SgrParams
SgrParams      := [0-9;]+
InlineOptional := '<optional>' Text '</optional>'
InlineAnyOf    := '<any-of>' InlineExpect+ '</any-of>'
InlineExpect   := '<expect>' Text '</expect>'
InlineRegex    := '<regex>' Pattern '</regex>'
Pattern        := /* Go regexp source, not containing unescaped '</regex>' */

BlockOptional  := '<optional>' newline Body '</optional>' newline
AnyOfBlock     := '<any-of>' newline ExpectBranch+ '</any-of>' newline
BlockRegex     := '<regex>' newline PatternLine '</regex>' newline
ContainsBlock  := '<contains>' newline ContainsBody '</contains>' newline
ContainsBody   := (ContainsFragmentLine | BlockStartWith | BlockEndWith)*
ContainsFragmentLine := PatternLine | LiteralLine
BlockStartWith := '<start-with>' newline ContainsBody '</start-with>' newline
BlockEndWith   := '<end-with>' newline ContainsBody '</end-with>' newline
ExpectBranch   := '<expect>' newline Body '</expect>' newline
Body           := Item*
Label          := [a-zA-Z][a-zA-Z0-9_]*
```

**Case sensitivity (§12 Q12):** all literal, hint, and regex matching is **case-sensitive** — no option in MVP.

**Only registered tags (§3) parse successfully.** Any other `<name>…</name>` form is a parse error.

### 5.4 Hint matching semantics

Hints are **documentation markers with literal matching**:

| Template fragment | Matches | Does not match |
|-------------------|---------|----------------|
| `<hint:id>abc-123</hint:id>` | `abc-123` only | `xyz`, `abc-124` |
| `id=<hint:id>abc</hint:id>` | `id=abc` only | `id=xyz` |
| `<hint:path>~/proj</hint:path>` | `~/proj` only | `/tmp/proj` |

**Rules:**

1. Text inside `<hint:label>…</hint:label>` must equal actual bytes at that position.
2. Open and close tag must use the same `label`.
3. On failure, report `hint:label` in the error (not "slot" or "placeholder").
4. For varying values, use `<any-of>`, `<optional>`, or `<contains>` — not hints.

### 5.5 ANSI color semantics

See §5.12 for rationale. Matcher compares inner text literally and independently verifies the ANSI open sequence for the specifier (named or raw `#…`).

### 5.6 Optional semantics (block + inline)

**Block optional:**

```
before
<optional>
  middle line
</optional>
after
```

| Rule | Detail |
|------|--------|
| Meta lines | `<optional>`, `</optional>` never matched against output |
| **Either / or** | Block is **absent** (zero lines consumed) **or** inner body is a **full literal match** — no partial consumption |
| Absent | Zero inner lines consumed from actual → pass |
| Present | Every inner line matched in order (literals, hints, inline optionals); all or nothing |
| Position | Between surrounding anchors; order matters |
| Nesting | Block optional inside block optional — **allowed** (§12 Q4) |
| **Adjacent blocks** | Back-to-back `<optional>` blocks are **separate** — **never merged** (§12 Q5) |

**Adjacent optionals (not merge):** each block keeps its own either-none-or-full-match semantics. Merging would change meaning (one combined block cannot skip its first inner lines).

```
before
<optional>
  line1
</optional>
<optional>
  line2
</optional>
after
```

| Actual between `before` and `after` | Match? |
|-------------------------------------|--------|
| (nothing) | yes — both absent |
| `line1` only | yes — first full, second absent |
| `line2` only | yes — first absent, second full |
| `line1` + `line2` | yes — both full |
| `line2` only if merged into one block | would be **no** — merge would wrongly forbid this |

**Inline optional:** only the wrapped `Text` may be absent; prefix/suffix on the same line remain literal (see §2.2).

### 5.7 Any-of semantics

```
<any-of>
<expect>
branch A line 1
branch A line 2
</expect>
<expect>
branch B line 1
</expect>
</any-of>
```

| Rule | Detail |
|------|--------|
| Meta lines | `<any-of>`, `</any-of>`, `<expect>`, `</expect>` never consume output |
| Matching | Try each branch in order; success if **one** branch fully matches the next actual lines |
| Consumption | Only the winning branch advances the actual-output cursor |
| Failure | Report **all branches** with per-branch deltas (§12 Q6) |
| Inner content | Each branch body follows normal `Item` rules (literals, hints, nested blocks) |

### 5.8 Contains semantics

```
<contains>
Usage: doctest
agent
<start-with>
  build
</start-with>
<end-with>
  skill
</end-with>
</contains>
```

| Rule | Detail |
|------|--------|
| Meta lines | `<contains>`, `</contains>` never consume output |
| Default fragment | Plain inner line → **full line** match: some line in actual must equal the fragment exactly |
| `<start-with>` | Wrapped text must match the **prefix** of some actual line (after block/inline disambiguation like `<optional>`) |
| `<end-with>` | Wrapped text must match the **suffix** of some actual line |
| Order | Fragments may match in any order relative to each other |
| Position | Block can appear anywhere in template; does not require contiguous actual region |
| Scope | `<start-with>` / `<end-with>` are valid **only** inside `<contains>` bodies |
| Combination | Often used as sole assertion in a template, or after exact-match header lines |

**Examples:**

| Fragment | Actual line | Match? |
|----------|-------------|--------|
| `agent` (default) | `  agent    Agent commands` | no — not full line |
| `  agent    Agent commands` | same | yes |
| `<start-with>agent</start-with>` | `  agent    Agent commands` | yes — prefix |
| `<end-with>management</end-with>` | `  skill    Skill management` | yes — suffix |

### 5.9 Regex semantics

**Syntax:** `<regex>PATTERN</regex>` — Go `regexp` (RE2). Pattern is **not** a wildcard tag; it is the regex source.

| Form | Behavior |
|------|----------|
| **Block** | `<regex>` / `</regex>` are meta lines; one inner pattern line **fully matches** exactly one actual line at cursor |
| **Inline** | Matched span in actual must fully match `PATTERN` (as if anchored to the span) |

**Examples:**

| Template | Actual | Match? |
|----------|--------|--------|
| `<regex>^\.+$</regex>` (whole line) | `..` | yes |
| same | `...` | yes |
| same | `.` then summary on same line | no — full line must match |
| `id=<regex>[0-9]+</regex>` | `id=42` | yes |

**Dot progress:** `<regex>^\.+$</regex>` replaces hand-rolled dot counting (§4.2).

### 5.10 Escaping (§12 Q16)

**Most `<` and `>` in CLI output do not need escaping.** Comparisons, redirection text, partial brackets, etc. are matched as literals unless the parser would read them as a **registered assert tag**.

Escape is required only when template text would otherwise parse as tag syntax (`<tag>`, `</tag>`, `<hint:…>`, `<optional>`, …).

**Syntax:** backslash before the angle bracket (template-side only; not expected in actual output):

| Template writes | Matches actual | Meaning |
|-----------------|----------------|---------|
| `2 < 3` | `2 < 3` | no escape — not tag-shaped |
| `\<optional>` | `<optional>` | literal text, not block optional |
| `see \<regex>` docs | `see <regex> docs` | literal `<regex>`, not regex tag |
| `\</optional>` | `</optional>` | literal close marker |

Rules:

1. `\` before `<` or `</` suppresses tag parsing for that marker.
2. Backslash is **not** part of actual output — template escape only.
3. Unregistered tag-shaped text without escape → parse error (not silently literal).
4. Actual output is never preprocessed; escaping is a template authoring concern.

### 5.11 Newlines

| Policy | Decision |
|--------|----------|
| Normalize `\r\n` → `\n` | **Always on** (§12 Q18) |
| Trailing newline | **Strict** — template and actual must agree; no automatic ignore (§12 Q9) |

### 5.12 ANSI color — `<ansi-color>` tag (preferred over `StripANSI`)

Declarative color assertions beat a global `StripANSI()` option:

| Approach | Problem |
|----------|---------|
| `StripANSI()` on entire actual | Strips color everywhere; cannot assert "this segment gray, that segment plain" |
| Raw escape bytes in template | Unreadable, brittle across terminals |
| **`<ansi-color gray>1 Cached</ansi-color>`** | Readable; asserts exactly which span has which color |

**Syntax:** `<ansi-color SPEC>text</ansi-color>` — inline only; inner `text` is literal.

**Matching:**

1. Find `text` at the current position in actual (same rules as literal/hint).
2. Resolve style tokens to expected **open SGR** sequence (left to right): `bold` → `\x1b[1m`; `gray` → `\x1b[90m`; `red` → `\x1b[31m`; `green` → `\x1b[32m`; `#90` → `\x1b[90m`; etc.
3. **Strict adjacency (§12 Q15):** open SGR bytes immediately precede `text`; exactly one `\x1b[0m` immediately follows `text`. No extra SGR between open sequence and text.
4. When color is disabled in actual (plain `text`), match **fails** — style assertion is explicit.

**Bold example:**

```
<ansi-color bold>Summary</ansi-color>
<ansi-color bold gray>1 Cached</ansi-color>
```

**Named example** (cached-gray test):

```
  (2 Run, 2 Pass, <ansi-color gray>1 Cached</ansi-color>, <ansi-color gray>0 Fail</ansi-color>)
```

**Raw example** (custom / future palette):

```
<ansi-color #38;5;208>CRITICAL</ansi-color>
```

Equivalent to asserting open `\x1b[38;5;208m` + literal `CRITICAL` + reset.

**Empty `<ansi-color red></ansi-color>`** — invalid (inner text required).

**Why both forms:** names stay readable for common doctest colors; `#` params assert uncommon sequences precisely without growing the reserved name list.

---

## 6. Match Modes

### 6.1 `MatchExact` (default — §12 Q8)

Entire actual output must equal the template (after `\r\n` normalization only). Best for `resp.Output`, `resp.Stdout` when output is bounded. This is the **default**; callers need not pass an option.

### 6.2 `MatchContains()` option (§12 Q17)

A **Go API option**, not a template tag. Entire template must appear as one **contiguous, order-preserving** substring of actual (after `\r\n` normalization).

```go
assert.Match(p, actual, assert.Contains())
```

Distinct from:

- **`MatchExact` (default):** actual equals template in full.
- **`<contains>` tag (§5.8):** listed fragments each appear **somewhere** in actual; **order-free**; fragments need not be adjacent.

### 6.3 Concrete comparison (Q17)

**Actual** (abbreviated noisy log):

```
[info] doctest build starting
[info] compiling pkg/a
[info] compiling pkg/b
..
  (2 Run, 2 Pass, 1 Cached, 0 Fail)
[info] done in 1.2s
```

**Template A** (exact summary block only):

```
..
  (2 Run, 2 Pass, 1 Cached, 0 Fail)
```

| Mode | Result | Why |
|------|--------|-----|
| `MatchExact` | **fail** | actual has extra lines before/after |
| `MatchContains()` | **pass** | template lines appear **contiguously** in the middle of the log |
| `<contains>` tag | **pass** | each fragment line found somewhere (order-free) |

**Template B** (fragments that are **not** adjacent in actual):

```
<contains>
compiling pkg/a
compiling pkg/b
</contains>
```

| Mode | Result | Why |
|------|--------|-----|
| `<contains>` tag | **pass** | each line found |
| `MatchContains()` with same two lines as template | **pass only if** those two lines are contiguous in actual; **fail** if another line sits between them |

**Template C** — why order-free `<contains>` can be looser:

Actual (lines reordered):

```
compiling pkg/b
compiling pkg/a
```

| Mode | Result |
|------|--------|
| `<contains>` tag with `pkg/a` + `pkg/b` fragments | **pass** (order-free) |
| `MatchContains()` with template `compiling pkg/a\ncompiling pkg/b` | **fail** (wrong order / not contiguous block matching template order) |

**When each is useful:**

| Tool | Use when |
|------|----------|
| `MatchExact` | Bounded stdout/stderr (typical doctest leaf) |
| `<contains>` tag | Help text, scattered keywords (`tests/help/…`) |
| `MatchContains()` | Long log; assert one **ordered multi-line excerpt** appears intact |

**MVP:** `MatchContains()` included (§12 Q17 = A). Default remains `MatchExact`.

---

## 7. API Design

```go
package assert

// Pattern is a parsed, immutable template.
type Pattern struct { /* ... */ }

// Parse parses template text into a Pattern.
func Parse(template string) (Pattern, error)

// MustParse parses or panics. For init/test literals.
func MustParse(template string) Pattern

// Match compares actual output against pattern. nil = success.
func Match(p Pattern, actual string, opts ...Option) error

// Output is a testing.T helper: parse + match + rich fatal.
func Output(t *testing.T, actual, template string, opts ...Option)

// Options
func Contains() Option            // MatchContains mode (contiguous subregion)
func NormalizeNewlines(v bool) Option  // \r\n → \n (default true); does not relax trailing newline
```

### 7.1 Error reporting

Failures must be actionable in doctest CI:

```
output mismatch at line 3:
  want: "$ cd <hint:path>~/proj</hint:path>"
  got:  "$ cd /tmp/proj"
              ^^^^^^^^
  hint:path expected "~/proj", got "/tmp/proj"
```

**Any-of failure (§12 Q6 — all branches):**

```
output mismatch at line 2: <any-of> — no branch matched
  actual: "error: connection reset"

  branch 1 (<expect>):
    want: "error: file not found"
    got:  "error: connection reset"
  branch 2 (<expect>):
    want: "error: permission denied"
    got:  "error: connection reset"
```

---

## 8. Data Model

No persistent storage. Pure in-memory AST + ephemeral match state.

### 8.1 Parsed template (immutable)

```go
type Pattern struct {
    Items []Item
}

type Item interface {
    isItem()
}

// One exact line (no inline tags).
type LiteralLine struct {
    Text string
}

// One line with interleaved literal/hint segments.
type PatternLine struct {
    Segments []Segment
}

type Segment interface {
    isSegment()
}

type Literal struct {
    Text string
}

type Hint struct {
    Label string // "id", "path", "duration"
    Text  string // literal text to match
}

type AnsiColor struct {
    Spec AnsiColorSpec // named or raw
    Text string        // literal inner text
}

type AnsiColorSpec struct {
    Tokens []string // "bold", "gray", or raw params "38;5;208" from # token
}

type RegexLine struct {
    Pattern string // Go regexp source
}

type InlineRegex struct {
    Pattern string
}

type ContainsFragment struct {
    Mode ContainsMatchMode // FullLine (default), StartWith, EndWith
    Text string
}

// Inline optional span within a PatternLine.
type InlineOptional struct {
    Text string // optional literal text
}

// Block optional: meta lines omitted; Items are match candidates.
type BlockOptional struct {
    Items []Item
}

// Any-of: one branch must match.
type AnyOfBlock struct {
    Branches []ExpectBranch
}

type ExpectBranch struct {
    Items []Item
}

// Order-free fragment assertions.
type ContainsBlock struct {
    Fragments []ContainsFragment
}
```

### 8.2 Match state (per call, ephemeral)

```go
type matchState struct {
    actual    string
    lines     []string
    cursor    int
    mode      matchMode // exact | contains
}
```

### 8.3 Storage layout

| Artifact | Location |
|----------|----------|
| Template source | String literals in Go, fenced blocks in ASSERT.md prose |
| Parsed pattern | In-memory; safe to cache/reuse |

---

## 9. Integration with Doctest

### 9.1 ASSERT.md usage (manual, MVP)

````markdown
## Expected Output

```output
..
  (2 Run, 2 Pass, 1 Cached, 0 Fail)
```

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
    assert.Output(t, resp.Output, `
..
  (2 Run, 2 Pass, 1 Cached, 0 Fail)
`)
}
```
````

### 9.2 Future: codegen from ` ```output ` blocks

`doctest agent fill-code` could generate `assert.Output(...)` from prose ` ```output ` fences. Not MVP.

### 9.3 Future: doctest skill doc

Add a section to `doc-spec` describing when to use output templates vs imperative asserts.

---

## 10. Test Plan

No implementation until design is confirmed. When approved, tests are mandatory.

### 10.1 Parsing (`assert_test.go` table-driven)

| ID | Input | Expected |
|----|-------|----------|
| P1 | `"hello\nworld"` | 2 literal lines |
| P2 | `"before\n<optional>\n</optional>\nafter"` | literal + block optional + literal |
| P3 | `"<hint:id>abc</hint:id>"` | pattern line, hint label `id` |
| P4 | `"id=<hint:id>abc</hint:id>"` | literal + hint segments |
| P5 | `"<id>abc</id>"` | parse error — bare tag not allowed |
| P6 | `"unclosed <hint:id>abc"` | parse error with position |
| P7 | `"Result: <optional>warn: </optional>OK"` | inline optional segment |
| P8 | `<any-of>` with two `<expect>` branches | `AnyOfBlock` with 2 branches |
| P9 | `<contains>` with 3 inner lines | `ContainsBlock` |
| P10 | `<hint:id>abc</hint:wrong>` | parse error — label mismatch |
| P11 | `<ansi-color gray>1 Cached</ansi-color>` | `AnsiColor` named `gray` |
| P12 | `<start-with>` inside `<contains>` | `ContainsFragment` mode prefix |
| P13 | `<ansi-color></ansi-color>` | parse error — empty inner |
| P14 | `<ansi-color #38;5;208>x</ansi-color>` | `AnsiColor` raw SGR `38;5;208` |
| P15 | `<ansi-color orange>x</ansi-color>` | parse error — unknown name; use `#` |
| P16 | `<regex>^\.+$</regex>` block + inline | `RegexLine` / `InlineRegex` |
| P17 | `status: <any-of><expect>ok</expect><expect>no</expect></any-of>` | inline any-of |
| P18 | `<ansi-color bold gray>x</ansi-color>` | tokens `["bold","gray"]` |
| P19 | `see \<optional> in docs` | literal `<optional>` in match text |

### 10.2 Matching — literals

| ID | Template | Actual | Result |
|----|----------|--------|--------|
| M1 | `OK` | `OK` | pass |
| M2 | `OK` | `OK\n` | fail — strict trailing newline (§5.10) |
| M2b | `OK\n` | `OK\n` | pass |
| M3 | `a\nb` | `a\nc` | fail line 2 |

### 10.3 Matching — optional (block + inline)

| ID | Template | Actual | Result |
|----|----------|--------|--------|
| O1 | `head\n<optional>\n</optional>\ntail` | `head\ntail` | pass — block absent |
| O2 | `head\n<optional>\nnoise\n</optional>\ntail` | `head\nnoise\ntail` | pass — block present |
| O3 | `head\n<optional>\n</optional>\ntail` | `headtail` | fail — missing newline |
| O4 | `Result: <optional>warn: </optional>OK` | `Result: OK` | pass — inline absent |
| O5 | `Result: <optional>warn: </optional>OK` | `Result: warn: OK` | pass — inline present |
| O6 | `Result: <optional>warn: </optional>OK` | `OK` | fail — literal prefix required |
| O7 | meta `<optional>` line alone | (any) | meta line never compared to actual |
| O8 | two adjacent block optionals | `line2` only between anchors | pass — separate semantics (§5.6) |
| O9 | one block optional inner | `line1` only when template has `line1` + `line2` | fail — partial inner not allowed |

### 10.4 Matching — any-of

| ID | Template | Actual | Result |
|----|----------|--------|--------|
| A1 | `<any-of><expect>a</expect><expect>b</expect></any-of>` | `a` | pass — first branch |
| A2 | same | `b` | pass — second branch |
| A3 | same | `c` | fail — no branch |
| A4 | `<any-of><expect>line1\nline2</expect><expect>alt</expect></any-of>` | two-line block matching branch 1 | pass |
| A5 | meta `<any-of>` / `<expect>` lines | (any) | meta lines never compared to actual |

### 10.5 Matching — hints

| ID | Template | Actual | Result |
|----|----------|--------|--------|
| H1 | `<hint:id>abc</hint:id>` | `abc` | pass |
| H2 | `<hint:id>abc</hint:id>` | `xyz` | fail — hint is literal, not wildcard |
| H3 | `id=<hint:id>abc</hint:id>` | `id=abc` | pass |
| H4 | `$ cd <hint:path>~/proj</hint:path>` | `$ cd /tmp/proj` | fail with `hint:path` in error |

### 10.6 Matching — contains

| ID | Template | Actual | Result |
|----|----------|--------|--------|
| C1 | `<contains>\na\nb\n</contains>` | `z\na\ny\nb` | pass — order-free full lines |
| C2 | same | `z\na` | fail — missing `b` |
| C3 | `<contains>\nfoo\n</contains>` | `x\nxfoo\n` | fail — `foo` is not a full line |
| C4 | `<contains>\n<start-with>foo</start-with>\n</contains>` | `x\nxfoo extra` | pass — prefix |
| C5 | `<contains>\n<end-with>bar</end-with>\n</contains>` | `line bar` | pass — suffix |
| C6 | meta `<contains>` line | (any) | meta line never compared |

### 10.6b Matching — ansi-color

| ID | Template | Actual | Result |
|----|----------|--------|--------|
| A1 | `<ansi-color gray>1 Cached</ansi-color>` | `\x1b[90m1 Cached\x1b[0m` | pass — named |
| A2 | same | plain `1 Cached` | fail — color expected |
| A3 | `<ansi-color green>2 Pass</ansi-color>` | `\x1b[32m2 Pass\x1b[0m` | pass — named |
| A4 | `<ansi-color #90>1 Cached</ansi-color>` | `\x1b[90m1 Cached\x1b[0m` | pass — raw equivalent to gray |
| A5 | `<ansi-color #38;5;208>warn</ansi-color>` | `\x1b[38;5;208mwarn\x1b[0m` | pass — raw 256-color |
| A6 | `<ansi-color #38;5;208>warn</ansi-color>` | `\x1b[90mwarn\x1b[0m` | fail — wrong SGR |
| A7 | `<ansi-color red></ansi-color>` | — | parse error — empty inner |
| A8 | `<ansi-color bold>Hi</ansi-color>` | `\x1b[1mHi\x1b[0m` | pass |
| A9 | `<ansi-color bold gray>1 Cached</ansi-color>` | `\x1b[1m\x1b[90m1 Cached\x1b[0m` | pass — strict order |
| A10 | `<ansi-color bold gray>1 Cached</ansi-color>` | `\x1b[90m1 Cached\x1b[0m` | fail — missing bold |

### 10.6c Matching — regex

| ID | Template | Actual | Result |
|----|----------|--------|--------|
| R1 | `<regex>^\.+$</regex>` | `..` | pass |
| R2 | same | `..` + summary glued | fail |
| R3 | `<optional><regex>^\.+$</regex></optional>` | summary only | pass — dots absent |

### 10.6d Matching — inline any-of

| ID | Template | Actual | Result |
|----|----------|--------|--------|
| X1 | `status: <any-of><expect>ok</expect><expect>err</expect></any-of>` | `status: ok` | pass |
| X2 | same | `status: err` | pass |
| X3 | same | `ok` | fail |

### 10.7 Modes & normalization

| ID | Scenario | Result |
|----|----------|--------|
| N1 | `Contains` — template embedded in long log | pass |
| N2 | `\r\n` actual, `\n` template | pass |
| N3 | strict newline on exact match | per M2/M2b |

### 10.8 Realistic doctest fixtures (post-MVP migration candidates)

| ID | Source test | Validates |
|----|-------------|-----------|
| R1 | `libdoc/build/tests/dot-progress/incremental` | multiline + summary line |
| R2 | `tests/help/top-level` | multi-line help |
| R3 | `libdoc/build/tests/output-color/color-enabled/cached-gray` | literal cached summary |
| R4 | `tests/help/top-level` | `<contains>` block |

### 10.9 How to run

```sh
go test ./assert/...
```

Doctest-style examples in `assert.go` doc comments for narrative cases; table tests for exhaustive edge cases.

---

## 11. Implementation Phases

### Phase 1 — MVP

- [ ] `Parse`, `MustParse`, `Match`, `Output`
- [ ] All tags in §3.1: `optional`, `any-of`, `expect`, `contains`, `regex`, `start-with`, `end-with`, `hint:label`, `ansi-color` (incl. `bold`, raw `#SGR`)
- [ ] Literal lines, pattern lines (inline hints + inline optional)
- [ ] Block meta lines do not consume actual output
- [ ] Hint literal matching; reject bare `<id>`-style tags
- [ ] `NormalizeNewlines` default on
- [ ] Line-numbered diff errors
- [ ] Strict trailing newline policy
- [ ] Selective `\<tag>` / `\</tag>` escaping (§5.10)
- [ ] `MatchContains()` option (contiguous subregion, §6.2)
- [ ] Table tests (§10.1–10.6b)

### Phase 2

- _(No open items — see §12.3 deferred topics.)_

### Phase 3

- [ ] ASSERT.md ` ```output ` codegen
- [ ] Migrate 2–3 existing ASSERT.md leaves as proof
- [ ] Doctest skill doc update

---

## 12. Open Questions

Sequential register of design decisions. **Needs confirm** = awaiting your answer before implementation.

### 12.1 Confirmed (Q1–Q19)

| # | Question | Decision |
|---|----------|----------|
| Q1 | Hint prefix — `<hint:label>` required? | **Yes** — bare `<id>` rejected |
| Q2 | Hint matching — wildcards? | **No** — inner text is literal only |
| Q3 | `<contains>` tag in MVP? | **Yes** |
| Q4 | Block `<optional>` nesting? | **Allowed** |
| Q5 | Adjacent block optionals — merge or separate? | **Separate, never merge**; each block is absent or full match |
| Q6 | `<any-of>` failure reporting? | **Show all branches** with want/got per branch |
| Q7 | Inline `<any-of>`? | **Yes** — MVP, same scoping as inline `<optional>` |
| Q8 | Default match mode? | **`MatchExact`** (full output equality after `\r\n` norm) |
| Q9 | Trailing newline policy? | **Strict** — no auto-ignore |
| Q10 | ANSI assertions mechanism? | **`<ansi-color SPEC>`** tag; named + raw `#SGR`; no global `StripANSI()` |
| Q11 | Dot progress variable dots? | **`<regex>`** tag (block + inline), e.g. `^\.+$` |
| Q12 | Case sensitivity? | **Strict** case-sensitive only |
| Q13 | `<contains>` fragment match mode? | **Default full line**; use `<start-with>` / `<end-with>` for prefix/suffix |
| Q14 | When to use `<hint:label>`? | Only when label aids readability; omit for obvious literals like `1 Cached` |
| Q15 | ANSI adjacency + bold? | **Strict** open SGR + reset; `bold` token; combine e.g. `<ansi-color bold gray>` |
| Q16 | Escaping `<` `>` in templates? | **Only when tag-shaped** — most literals need no escape; use `\<tag>` / `\</tag>` to match literal angle markers (§5.10) |
| Q17 | `MatchContains()` Go option (§6.3)? | **Yes — MVP** (contiguous ordered subregion; default stays `MatchExact`) |
| Q18 | `\r\n` → `\n` normalization? | **Always on** |
| Q19 | Start implementation now? | **No** — stay in doc phase (§14) |

### 12.2 Needs your confirm

_None — Q1–Q19 all confirmed. Implementation blocked only by §14 (explicit **go ahead**)._

### 12.3 Deferred by design (not blocking MVP — no confirm required unless you disagree)

| Topic | Current plan |
|-------|----------------|
| ASSERT.md ` ```output ` codegen | Phase 3 |
| Migrate existing ASSERT.md leaves | Phase 3 |
| Doctest skill doc update | Phase 3 |
| Block `<regex>` with multiple inner pattern lines | MVP: **one** pattern line per block; extend later if needed |
| Additional ANSI color names beyond `red`/`green`/`gray`/`bold` | Use raw `#SGR` params |

---

## 13. Non-Goals (for now)

- Bare tags (`<id>`, `<cached>`, `<run>`) and user-defined tag names
- Wildcard / variable slots inside hints
- Global `StripANSI()` preprocessing (use `<ansi-color>` per segment instead)
- Free-form regex outside `<regex>` tags
- Snapshot file I/O and auto-update
- Structural matching of JSON/XML inside output
- Replacing all non-output assertions (exit code, error types, helper predicates)

---

## 14. Approval Gate

Implementation starts only after explicit **"go ahead"** on a stabilized version of this document.

On approval, the implementer will:

1. Build `assert/` per §7–8
2. Run tests per §10
3. Optionally migrate one ASSERT.md leaf as a reference

---

## Changelog

| Date | Change |
|------|--------|
| 2026-06-26 | Initial draft from brainstorm |
| 2026-06-26 | Replace `<>` with `<optional>`; add `<any-of>`/`<expect>`; block meta vs inline scoping |
| 2026-06-26 | §2.5 tag taxonomy; clarify placeholders vs meta; reject `<dots>` anti-pattern |
| 2026-06-26 | §3 Included Tags; `<hint:label>` literal hints; `<contains>`; reject bare `<id>`/`<cached>` |
| 2026-06-26 | Q4/Q8/Q9/Q10/Q13 decisions; `<ansi-color>`; `<start-with>`/`<end-with>`; strict newlines |
| 2026-06-26 | `<ansi-color #SGR>` raw codes alongside named colors |
| 2026-06-26 | Q5: adjacent block optionals stay separate (none or full match; no merge) |
| 2026-06-26 | Q6/Q7/Q11/Q12/Q15; `<regex>` block+inline; inline `<any-of>`; `bold` in `<ansi-color>` |
| 2026-06-26 | §12 reordered Q1–Q19; split confirmed vs needs confirm |
| 2026-06-26 | Q16 selective `\<tag>` escape; Q18/Q19 confirmed; §6.3 Q17 concrete examples |
| 2026-06-26 | Q17 confirmed — `MatchContains()` in MVP |
| 2026-06-26 | Status → Implemented; §16 doc/skill advocacy plan (awaiting go ahead) |
| 2026-06-26 | §16 implemented: output-assert skill, design/code/doc specs, review skill |

---

## 16. Doc & skill advocacy plan

> **Status:** **Done** — spec/skill advocacy shipped (§16.2).

### 16.1 Locked decisions (A–E)

| ID | Decision |
|----|----------|
| **A** | **Full tag spec** embedded in design spec (§3 content), not link-only |
| **B** | `## Expected Output` in ASSERT.md — **advertised, not mandatory** (`doctest vet` unchanged) |
| **C** | New standalone skill: **`doctest skill output-assert show\|install`** |
| **D** | Review flags legacy `strings.Contains` / hand-rolled parsing as **suggestion** severity, not must-fix |
| **E** | This file: status **Implemented** + authoring quick start (above) |

### 16.2 Files to edit (on go ahead)

| File | Change |
|------|--------|
| `doc/DOCTEST_OUTPUT_ASSERT.md` | **New** — skill-facing copy of §3 + §5–§7 + authoring (from this doc); YAML frontmatter `name: doctest-output-assert` |
| `doc/doc.go` | `//go:embed DOCTEST_OUTPUT_ASSERT.md`; `Content()` case |
| `libdoc/spec/spec.go` | `entries["output-assert"]` |
| `libdoc/cli/cli.go` | List `output-assert` in skill help + usage |
| `doc/snippets/DOCTEST_DESIGN_SPEC.md` | New **## Output assertions** — full §3 tag registry + when-to-use table + `assert.Output` example; pointer to `doctest skill output-assert show` |
| `doc/DOC_STYLE_TEST_CODE_SPECIFICATION.md` | Replace `strings.Contains` exemplars; document `import assert`, `Output`, `Match`+`Contains()` |
| `doc/DOC_STYLE_TEST_SPECIFICATION.md` | ASSERT format: optional `## Expected Output` (recommended prose mirror) |
| `doc/DOCTEST_REVIEW.md` | Checklist + anti-patterns + report item **Output assertions**; suggestion severity for legacy patterns |

**No changes:** `doctest vet` rules, existing ASSERT.md leaves (unless user asks migration later).

### 16.3 New skill shape (`DOCTEST_OUTPUT_ASSERT.md`)

```yaml
---
name: doctest-output-assert
description: >-
  Output assert DSL for doctest ASSERT.md — template tags, API, and migration
  from strings.Contains.
---
```

Body sections (condensed from this design doc):

1. When to use assert vs literals
2. **Full §3 Included Tags** (copy verbatim)
3. API: `Parse`, `Match`, `Output`, `Contains()` option
4. ASSERT.md pattern (`## Expected Output` optional + `func Assert`)
5. Migration table: `strings.Contains` loop → `<contains>`; dots → `<regex>`; platform errors → `<any-of>`; ANSI helpers → `<ansi-color>`
6. `go test ./assert/...` / `doctest test ./assert/tests/...` for examples

### 16.4 Review skill additions (sketch)

**Checklist (new):**

- [ ] Structured CLI output uses `github.com/xhd2015/doctest/assert`
- [ ] Non-trivial templates mirrored in `## Expected Output` when present
- [ ] Legacy `strings.Contains` / `Index`/`Count` on output fields → **suggestion** to migrate

**Anti-patterns (suggestion):**

| Legacy | Prefer |
|--------|--------|
| `strings.Contains` loop | `<contains>` + `<start-with>` |
| Dot/summary parsing | `<regex>^\.+$</regex>` + literals |
| Dual platform `Contains` | `<any-of>` |
| `metricIsColored` | `<ansi-color bold gray>` |

**Report section:** add **Output assertions** per root (between scenario quality and mechanical checks).

### 16.5 CLI wiring (on go ahead)

```
doctest skill --list          # includes output-assert
doctest skill output-assert show
doctest skill output-assert install --cursor
```

Update `tests/skill/list/ASSERT.md` want list if it enumerates skills (mechanical, user said no new tests — **only if list test breaks**).
