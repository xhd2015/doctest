## Code blocks

`SETUP.md` and `ASSERT.md` may contain ```go...``` blocks.

## DOCTEST.md

Root of a decision-tree of scenarios. Root Go block defines the shared contract.

### DSN (Domain Specific Notion)

Every root must have `# DSN (Domain Specific Notion)` (`doctest vet` rejects missing).

**DSN is a domain sketch, not a DSL and not the test plan.** Scannable mental
model of who acts and what they do — selective (load-bearing actors only);
leaves/asserts carry precision.

| Dialect | Where | Shape |
|---------|--------|--------|
| **Root DSN sketch** | `# DSN …` in root `DOCTEST.md` | **Participants** + **Behaviors** prose (no ``` fences) |
| **Scenario sketch** | `# Scenario` ``` block in each `SETUP.md` | One path: `A -> B -> effect` + optional `#` captions |

**Root sketch:** optional one-line purpose; ~3–8 bold-named participants; behavior
bullets (light `A -> B` OK). No code fences, asserts, or MECE tree plan.
About right: one screen. Too thin: rename dirs. Too heavy: every edge / code tour.
Scenario shape: see **Scenario**.

### Contracts (define once)

| Symbol | Where | Rule |
|--------|-------|------|
| `type Request`, `type Response`, `func Run` | Root `DOCTEST.md` Go only | Never redefine in child SETUP |
| `func Setup` | Root and/or grouping/leaf `SETUP.md` | Root→leaf chain; body not stub `return nil` |
| Helpers | Any `SETUP.md` Go block | Inherited; child must not redefine same name |

Root `SETUP.md` must not redefine Request/Response/Run. Generated signatures
include `d *session.Doctest` (see **Parallel-safe harness**).

### Nested DOCTEST.md

Subdir with its own `DOCTEST.md` = **self-contained root** — no inheritance
crosses (types, helpers, Setup, Run). Own contract + Setup chain + binary
resolution. Parent Setup never runs under a nested root.

**When:** groups cannot share the same `Run(Request, Response)` → separate trees.

`d.DOCTEST_ROOT` is the nested root (join `".."` for module root).
`d.DOCTEST_SESSION_ID` still shared for one `doctest test` run.

### Test-first design & layers

Trees/`Run` target **testable surfaces** (easy to test, not e2e-only glue).

| Layer | Default use | Notes |
|-------|-------------|--------|
| **L1** | Pure / flat edges | `*_test.go` tables |
| **L2** ★ | Multi-factor APIs + short CLI | Library or `cli.RunWithWriter` |
| **L3** | Full integration only | Separate process; `label: e2e` |

Default **L2**. Short path → never L3. Deep dive:
`doctest skill design-principle --show`.

### Tree Organization

**MECE** (Mutually Exclusive, Collectively Exhaustive):

1. **Mutually exclusive siblings** — distinct branches, no duplicate coverage.
2. **Collectively exhaustive (pragmatic)** — all **meaningful** outcomes; skip trivia.

**Significance:** high-impact factors high in the tree; minor variants lower.
Parent→child narrows `Request` by one or a few params (most significant first).
Siblings MECE for the factor split at that level.

## SETUP.md

### Scenario

First section must be `# Scenario` (vet). Tag line: `**Feature**: …` or
`**Bug**: …`, then a fenced **scenario sketch** — one path through the root DSN
(not a full re-inventory). Prefer `A -> B -> effect` / `<-`; `# comment` above
unclear hops; only root participants; ~2–8 lines (not empty / Feature-only /
whole DSN paste).

<example-of-SETUP.md>
# Scenario

**Feature**: agent commands use fake Codex instead of a real LLM

```
# agent reads requirement, invokes Fake Codex, writes output
doctest agent <cmd> --requirement req.md -> Fake Codex -> generated code

# session state tracked in event files
doctest <- Fake Codex (session id, events, progress)
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error { ... }
```
</example-of-SETUP.md>

### Code

`SETUP.md` ends with a Go block. Non-root must define `Setup`; root `DOCTEST.md`
defines `Run`. See **Contracts**.

# Inheritance

**Setup chain:** root → intermediates → leaf on the same `*Request`. **Run** =
root only for that tree (else **Nested**). **Helpers:** root = tree-wide;
grouping = subtree; no same-name redefinition. Types/Run: **Contracts**.

## Parallel-safe harness

Leaves use `t.Parallel()` in one process.

**Forbidden** in harness and L2 product: (1) unprotected shared globals /
reassigning `os.Stdout|Stderr|Stdin`; (2) process env/cwd —
`os.Setenv`/`Unsetenv`, `os.Chdir`, `t.Setenv`, `t.Chdir`, `syscall.Setenv`.
**Prefer** inject opts / `req` / child `cmd.Env`·`Dir`. Detail:
`doctest skill lint --show`, `doctest skill review --show` (Common gotchas).

**Context only via `d *session.Doctest`:** `d.DOCTEST_SESSION_ID`,
`d.DOCTEST_ROOT`, `d.DOCTEST_CASE`. No free inject vars; no author getenv of
those names. `doctest vet` rejects env reads. Pass `d` into path/session
helpers. Process cwd **undetermined** — absolute paths from `d`.

## Session-scoped shared setup (`d.DOCTEST_SESSION_ID`)

Amortize expensive setup **once per `doctest test`** with a cache keyed by
`d.DOCTEST_SESSION_ID`: one cache dir; file lock; ready marker; share only safe
artifacts (per-leaf temp for mutated state). Document layout in root SETUP
**Preconditions**.

```go
func sessionCacheDir(d *session.Doctest) string {
    return filepath.Join(os.TempDir(), "my-feature-doctest-"+d.DOCTEST_SESSION_ID)
}
// buildOnce: flock + ready marker; return artifact path under cacheDir
```

## ASSERT.md

Optional YAML frontmatter for run profile:

```yaml
---
label: ui-automation, slow
explanation: AX tree poll; compile and link ~25s
---
```

Every leaf needs:

```
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error)
```

Fail with `t.Fatal`/`t.Fatalf`. Unexported APIs: `TestExported_` wrappers.

## Output assertions

CLI/text (`resp.Stdout`, `Stderr`, `Output`, `Summary`, …): prefer
`github.com/xhd2015/doctest/assert` over `strings.Contains` loops.

```sh
doctest skill output-assert --show    # full registry + API
```

| Situation | Approach |
|-----------|----------|
| Bounded stdout/stderr | `assert.Output(t, actual, template)` |
| Variable port/path/user | `__PLACEHOLDER__` in YAML header |
| Skippable middle | `...N lines omitted...` |
| Flexible / platform line | raw regexp or `(linux\|darwin)` |
| ANSI segments | `<ansi-color bold gray>…</ansi-color>` |

Prefer **v3** (YAML header; each line raw Go regexp; strict line-by-line). No
`<contains>` / `assert.Contains()` for new tests. Optional `## Expected Output`
prose mirror before `## Expected`. Full constructs: `output-assert --show`.

```go
assert.Output(t, resp.Stdout, `---
version: 3
---
Usage: mytool
  build
  test
`)
```

## Test Fixture Data

Fixtures as files (single beside `ASSERT.md`; multiple in `testdata/`). Read via
`d` (e.g. `filepath.Join(d.DOCTEST_CASE, "input.txt")`). See **Parallel-safe
harness** and **Test-first design & layers**.

> Full spec: `doctest skill doc-spec --show` && `doctest skill code-spec --show` &&
> `doctest skill design-principle --show` && `doctest skill output-assert --show`
