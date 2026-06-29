## Code blocks

Both `SETUP.md` and `ASSERT.md` may contain ```go...``` go code blocks. 

## DOCTEST.md

`DOCTEST.md` marks the root of the test tree. The whole test tree rooted from where `DOCTEST.md` begins forms a large decision tree. The root `DOCTEST.md` Go block (final content) must define `type Request`, `type Response`, and `func Run` — these are shared by all descendants and must not be redefined. Root `SETUP.md` (when present) must not redefine them.

### DSN (Domain Specific Notion)

Every root `DOCTEST.md` must include a `# DSN (Domain Specific Notion)` section.
`doctest vet` rejects roots missing this section.
DSN is like a DSL, but less formal — it models the target under test as a
normal human mental model. It defines **participants** (actors, components,
subsystems) and their **behaviors** (what each participant does, how they
interact), written as plain prose (no code blocks). Think of it as a prose
sketch of the domain that helps readers understand what the test tree is
exercising.

Each `SETUP.md`'s `# Scenario` section (see below) wraps a snippet of this
DSN model in a ``` block, showing the subset of participants and behaviors
relevant to that particular scenario.

### Nested DOCTEST.md

A subdirectory that contains its own `DOCTEST.md` becomes a **self-contained test root**. The doctest runner stops walking at `DOCTEST.md` boundaries and treats each root independently — **no inheritance crosses a `DOCTEST.md` boundary**.

A nested root must be entirely self-sufficient:
- Its `DOCTEST.md` Go block must define its own `Request`/`Response` types and `func Run`
- It must provide `Setup` or let descendant SETUPs provide `Setup`
- Any external binaries (e.g., the doctest binary for `req.Bin`) must be built or resolved within that root's own `Setup`
- The parent root's `Setup` is never executed for leaves under a nested DOCTEST.md
- Paths like `DOCTEST_ROOT/..` shift — from a deeper root, use `DOCTEST_ROOT/../..` to reach the module root
- `DOCTEST_SESSION_ID` is shared within one `doctest test` run

### When to Create a Nested DOCTEST.md

If two test groups cannot share the same `Run(Request, Response)` contract,
they must be separate test trees — each rooted at its own `DOCTEST.md`. This
happens when different scenarios call different functions, services, or
execution strategies.

### Tree Organization

Doctest trees are decision trees. Design them using the **MECE** principle
(Mutually Exclusive, Collectively Exhaustive):

1. **Mutually exclusive siblings**: each sibling dir tests a distinct branch —
   no two siblings should cover the same scenario.
2. **Collectively exhaustive siblings (pragmatic)**: at each split, sibling
   dirs should cover all **meaningful** outcomes for that factor. Do not force
   branches for trivial, duplicate, or low-value cases; avoid gaps where an
   important case has no branch.

**Significance ordering**: place factors with the **largest impact on behavior
or outcome** at **higher** ancestor dirs; reserve **least significant** factors
(minor variants, edge values) for **lower** descendant dirs. Each parent→child
step should narrow `Request` by one or a few params, preferring the
highest-impact unresolved factor first.

3. **Parent → child dirs**: scenarios become more concrete by narrowing one or
   a few params from `Request` (most significant first).
4. **Sibling dirs**: must be MECE — mutually exclusive and pragmatically
   collectively exhaustive for the factor being split at that level.

## SETUP.md

### Scenario

Every `SETUP.md` must include a `# Scenario` section as its **first** section.
`doctest vet` rejects any `SETUP.md` that does not start with `# Scenario`.
This section starts with a tag line — either `**Feature**: <description>` or `**Bug**: <description>` — followed by a ``` block containing a DSN snippet
(from the root `DOCTEST.md`'s DSN model) that sketches the mental model with
annotated pipeline lines (`# comment` above each `->` / `<-` line).

<example-of-SETUP.md>
# Scenario

**Feature**: agent commands use fake Codex instead of a real LLM

```
# agent reads requirement, invokes Fake Codex, writes output
doctest agent <cmd> --requirement req.md -> Fake Codex -> generated code

# session state tracked in event files
doctest <- Fake Codex (session id, events, progress)
```

## Preconditions
- Agent commands must be able to use fake Codex instead of a real LLM.

## Steps
1. Lookup `fake-codex` from PATH; skip if not installed.

## Context
- ...

```go
func Setup(t *testing.T, req *Request) error { ... }
```
</example-of-SETUP.md>

### Code

Every `SETUP.md` must have a Go block as **final content**. Child must not redefine `Request`/`Response`/`Run`.

| Function | Signature | Notes |
|----------|-----------|-------|
| `Setup` | `(t *testing.T, req *Request) error` | Called root→leaf before `Run`; body must not be stub |
| `Run` | `(t *testing.T, req *Request) (*Response, error)` | **DOCTEST.md only**, cannot be redefined by descendants |

Root `DOCTEST.md` must define `Run`. Non-root `SETUP.md` must define `Setup`. Signatures must match exactly. `func Setup` body must not be a stub (`return nil`).

# Inheritance

## Setup Chain

`Setup` functions accumulate from root → leaf: every ancestor `SETUP.md`'s `func Setup` is called in order (root first, then each intermediate directory, then the leaf). Each `Setup` receives the same `*Request` and can modify it incrementally.

## Run Resolution

`func Run` is defined **only in the root** `DOCTEST.md` Go block. All leaves share the same `Run` function. Descendants must not redefine it.

If two test scenarios need a different `Run`, they require separate test trees, each rooted at its own `DOCTEST.md`.

## Type and Helper Scoping

- `type Request`, `type Response`, and `func Run` are defined once in the root `DOCTEST.md` Go block. Child `SETUP.md` files **must not** redefine them.
- Besides `func Setup`, each `SETUP.md` Go block may declare **helper functions** for its subtree:
  - **Root `SETUP.md`** — helpers shared by every test in the tree (e.g. build binary, write fixtures).
  - **Grouping `SETUP.md`** — helpers shared only by that node's descendants.
- Helpers in ancestor `SETUP.md` files are inherited by descendants. A child **must not** redefine a helper with the same name.

## DOCTEST.md Boundary

A `DOCTEST.md` file creates an **inheritance firewall**. No code, types, helpers, `Run`, or `Setup` functions cross a `DOCTEST.md` boundary. Each tree rooted at a `DOCTEST.md` is a self-contained decision tree with its own `Request`/`Response`/`Run` and its own setup chain.

## ASSERT.md

### Run profile labels (optional YAML frontmatter)

Runnable leaves may prefix `ASSERT.md` with YAML frontmatter:

```yaml
---
label: ui-automation, slow
explanation: AX tree poll; compile and link ~25s
---
```

### Assert function

Every `ASSERT.md` must have a `func Assert`. Signature must match exactly:

```
func Assert(t *testing.T, req *Request, resp *Response, err error)
```

Fail via `t.Fatal`/`t.Fatalf`.

Import target package directly. For unexported functions, use **`TestExported_`** prefix:
`func TestExported_foo() { foo() }` — then `import "mypkg"; mypkg.TestExported_foo()` in the code block.

## Output assertions

When a leaf asserts **CLI or text output** (`resp.Stdout`, `resp.Stderr`, `resp.Output`, `resp.Summary`, …), prefer the **`github.com/xhd2015/doctest/assert`** template DSL over `strings.Contains` loops and hand-rolled parsing.

```sh
doctest skill output-assert show    # full tag registry + API
```

### When to use

| Situation | Approach |
|-----------|----------|
| Bounded stdout/stderr | `assert.Output(t, actual, template)` |
| Bounded output with scattered required lines | ``assert.Output(t, actual, `<contains>...</contains>`)`` |
| Excerpt in a long log | `assert.Match(p, actual, assert.Contains())` |
| Scattered help keywords | `<contains>` + `<start-with>` |
| Platform-specific errors | `<any-of><expect>…</expect></any-of>` |
| Variable dot lines | `<regex>^\.+$</regex>` |
| ANSI-colored segments | `<ansi-color bold gray>…</ansi-color>` |

### Recommended prose mirror (not required by vet)

Optional `## Expected Output` section with a fenced template block before `## Expected` — aids review and readability.

### ASSERT example

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

### Included tags (full registry)

| Tag | Form | Consumes actual? | Matching |
|-----|------|------------------|----------|
| `<optional>` | block or inline | block meta: no | absent or full inner match |
| `<any-of>` | block or inline | block meta: no | one `<expect>` branch |
| `<expect>` | inside `<any-of>` | block meta: no | branch delimiter |
| `<regex>` | block or inline | yes | Go regexp full match |
| `<contains>` | block only | no (meta) | fragments anywhere; default full line |
| `<start-with>` | inside `<contains>` | no | line prefix |
| `<end-with>` | inside `<contains>` | no | line suffix |
| `<hint:label>` | inline | yes | literal match; label for docs |
| `<ansi-color>` | inline | yes | literal text + strict ANSI |

**Block meta** (standalone lines): `<optional>`, `<any-of>`, `<expect>`, `<contains>`, `<regex>` (+ closers). Meta lines never consume output except `<regex>` inner pattern line.

**Inline:** `<optional>…</optional>`, `<hint:label>…</hint:label>`, `<ansi-color SPEC>…</ansi-color>`, `<any-of>…</any-of>`, `<regex>…</regex>`, `<start-with>` / `<end-with>` inside `<contains>`.

Ordinary lines inside `<contains>` may use inline pattern tags such as
`<any-of>`, `<optional>`, `<regex>`, and `<hint:...>`.

Use `assert.Output(t, actual, template)` for templates whose top-level form is
`<contains>`. Reserve `assert.Match(p, actual, assert.Contains())` for finding a
contiguous excerpt in a larger output, and do not combine it with a top-level
`<contains>` block.

Matcher DSL tags are test syntax only. They must not be required in actual
product stdout/stderr unless the feature explicitly defines those strings as
user-facing output.

**`<ansi-color>` tokens:** `bold`, `red`, `green`, `gray`, or raw `#SGR` (e.g. `#90`, `#38;5;208`). Combined left-to-right: `<ansi-color bold gray>1 Cached</ansi-color>`.

**Rejected:** bare `<id>`, `<cached>`, `<run>`, user-defined tags, `<>` syntax.

**Escaping:** only tag-shaped text needs `\<tag>` / `\</tag>` in templates; plain `2 < 3` needs no escape.

**Avoid in new tests:** `strings.Contains` loops for structured output; `strings.Index`/`Count` for dots/summary; ad-hoc ANSI helpers when `<ansi-color>` applies.

## Test Fixture Data

Abstract fixture data into standalone files, not inline code.

- Single file → place alongside `ASSERT.md`
- Multiple files → place in `testdata/` alongside `ASSERT.md`

Code reads them with directly filename reference as each `ASSERT.md` runs in its own directory.

> Full spec, run: `doctest skill doc-spec show` && `doctest skill code-spec show` && `doctest skill output-assert show`
