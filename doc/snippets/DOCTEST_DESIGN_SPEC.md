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
- `DOCTEST_SESSION_ID` (injected variable) is shared within one `doctest test` run

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

## Session-scoped shared setup (`DOCTEST_SESSION_ID`)

`doctest test` runs leaf packages in parallel. When many leaves repeat the same
expensive setup (building binaries, creating seed archives, downloading fixtures),
amortize it **once per `doctest test` invocation** with a file-based cache keyed
by `DOCTEST_SESSION_ID`.

### Rules

1. **`DOCTEST_SESSION_ID` is an injected variable** — doctest defines it in every
   generated test. Reference it directly in harness helpers; do **not** call
   `os.Getenv("DOCTEST_SESSION_ID")`. `doctest vet` rejects env reads.
2. **One cache dir per run** — e.g. `$TMPDIR/my-feature-doctest-<DOCTEST_SESSION_ID>/`.
   All leaves in the invocation share the same session id.
3. **File lock for first-time population** — use `syscall.Flock` on a lock file
   inside the cache dir so parallel packages serialize the first build/seed; later
   packages wait, then reuse artifacts.
4. **Ready markers** — after populating, write a sentinel (e.g. `binaries.ready`)
   so waiters skip rebuild when artifacts exist.
5. **Per-leaf isolation unchanged** — share only artifacts safe across leaves
   (compiled binaries, generic seed archives). Keep per-leaf temp dirs for
   mutated state (server home, agent home, custom excludes/includes).

### Pattern (root `SETUP.md` helpers)

```go
func sessionCacheDir() string {
    return filepath.Join(os.TempDir(), "my-feature-doctest-"+DOCTEST_SESSION_ID)
}

func withFileLock(t *testing.T, lockPath string, fn func() error) error {
    // open lockPath, syscall.Flock(LOCK_EX), defer LOCK_UN, run fn
}

func buildOnce(t *testing.T, cacheDir string) string {
    lock := filepath.Join(cacheDir, "build.lock")
    ready := filepath.Join(cacheDir, "binaries.ready")
    bin := filepath.Join(cacheDir, "my-tool")
    withFileLock(t, lock, func() error {
        if fileExists(ready) && fileExists(bin) {
            return nil
        }
        // go build -o bin ...
        return os.WriteFile(ready, []byte("ok"), 0644)
    })
    return bin
}
```

Call `buildOnce` (or `ensureSessionDefaultArchive`) from root `Run` or shared
helpers. Leaves that need a custom seed still build their own archive; leaves
that only need the default seed reuse the session archive path.

Document the cache layout in root `SETUP.md` **Preconditions** so reviewers see
what is shared vs per-leaf.

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
| Bounded stdout/stderr | `assert.Output(t, actual, template)` — strict full match |
| Variable port, path, user | `__PLACEHOLDER__` in YAML header |
| Skippable middle section | `...N lines omitted...` |
| Flexible single line | regex line (e.g. `^\.+$`) |
| Platform-specific one-liner | regex alternation `(linux\|darwin)` |
| ANSI-colored segments | `<ansi-color bold gray>…</ansi-color>` |

v2 templates start with `version: 2` YAML header. Strict line-by-line match only — no `<contains>` or `assert.Contains()`.

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
    assert.Output(t, resp.Stdout, `---
version: 2
---
Usage: mytool
  build
  test`)
}
```

### v2 template constructs

| Construct | Form | Matching |
|-----------|------|----------|
| Placeholder | `__NAME__` in header + body | `type=string` or `type=number` |
| Omit marker | `...N lines omitted...` | skip exactly N actual lines |
| Regex line | line with regex-intent signals | Go regexp full line match |
| Pattern line | default literal line | literal + placeholders + color |
| `<ansi-color>` | inline only | strict ANSI wrap (same tokens as v1) |

Placeholder header: `__PORT__: type=number, example=8901, a port` — `k=v` metadata, trailing text is human explanation.

**`<ansi-color>` tokens:** `bold`, `red`, `green`, `gray`, or raw `#SGR` (e.g. `#90`, `#38;5;208`). Combined left-to-right: `<ansi-color bold gray>1 Cached</ansi-color>`.

v1 tag templates (`<contains>`, `<any-of>`, …) still parse via `legacy_v1` when no `version: 2` header. Prefer v2 for new tests.

**Avoid in new tests:** `strings.Contains` loops for structured output; `strings.Index`/`Count` for dots/summary; ad-hoc ANSI helpers when `<ansi-color>` applies.

## Test Fixture Data

Abstract fixture data into standalone files, not inline code.

- Single file → place alongside `ASSERT.md`
- Multiple files → place in `testdata/` alongside `ASSERT.md`

Code reads them with directly filename reference as each `ASSERT.md` runs in its own directory.

> Full spec, run: `doctest skill doc-spec show` && `doctest skill code-spec show` && `doctest skill output-assert show`
