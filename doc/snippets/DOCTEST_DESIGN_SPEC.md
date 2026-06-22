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

1. **Parent → child dirs**: scenarios become more concrete by narrowing one or a few params from `Request`.
2. **Sibling dirs**: must be mutually exclusive — each tests a different scenario branch.

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

Every `ASSERT.md` must have a `func Assert`. Signature must match exactly:

```
func Assert(t *testing.T, req *Request, resp *Response, err error)
```

Fail via `t.Fatal`/`t.Fatalf`.

Import target package directly. For unexported functions, use **`TestExported_`** prefix:
`func TestExported_foo() { foo() }` — then `import "mypkg"; mypkg.TestExported_foo()` in the code block.

## Test Fixture Data

Abstract fixture data into standalone files, not inline code.

- Single file → place alongside `ASSERT.md`
- Multiple files → place in `testdata/` alongside `ASSERT.md`

Code reads them with directly filename reference as each `ASSERT.md` runs in its own directory.

> Full spec, run: `doctest skill doc-spec show` && `doctest skill code-spec show`