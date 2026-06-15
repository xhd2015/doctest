## Code blocks

Both `SETUP.md` and `ASSERT.md` may contain ```go...``` go code blocks. 

## DOCTEST.md

`DOCTEST.md` marks the root of the test tree. The whole test tree rooted from where `DOCTEST.md` begins forms a large decision tree. The root `SETUP.md` must define `type Request`, `type Response`, and `func Run` — these are shared by all descendants and must not be redefined.

### Nested DOCTEST.md

A subdirectory that contains its own `DOCTEST.md` becomes a **self-contained test root**. The doctest runner stops walking at `DOCTEST.md` boundaries and treats each root independently — **no inheritance crosses a `DOCTEST.md` boundary**.

A nested root's `SETUP.md` must be entirely self-sufficient:
- It must define its own `Request`/`Response` types and `func Run`
- It must provide `Setup` or let descendant SETUPs provide `Setup`
- Any external binaries (e.g., the doctest binary for `req.Bin`) must be built or resolved within that root's own `Setup`
- The parent root's `Setup` is never executed for leaves under a nested DOCTEST.md
- Paths like `DOCTEST_ROOT/..` shift — from a deeper root, use `DOCTEST_ROOT/../..` to reach the module root

### When to Create a Nested DOCTEST.md

If two test groups cannot share the same `Run(Request, Response)` contract,
they must be separate test trees — each rooted at its own `DOCTEST.md`. This
happens when different scenarios call different functions, services, or
execution strategies.

### Tree Organization

1. **Parent → child dirs**: scenarios become more concrete by narrowing one or a few params from `Request`.
2. **Sibling dirs**: must be mutually exclusive — each tests a different scenario branch.

## SETUP.md

Every `SETUP.md` must have a Go block as **final content**. Child must not redefine `Request`/`Response`/`Run`.

| Function | Signature | Notes |
|----------|-----------|-------|
| `Setup` | `(t *testing.T, req *Request) error` | Called root→leaf before `Run`; body must not be stub |
| `Run` | `(t *testing.T, req *Request) (*Response, error)` | **Root only**, cannot be redefined by descendants |

Root must define `Run`. Non-root SETUP.md must define `Setup`. Signatures must match exactly. `func Setup` body must not be a stub (`return nil`).

# Inheritance

## Setup Chain

`Setup` functions accumulate from root → leaf: every ancestor `SETUP.md`'s `func Setup` is called in order (root first, then each intermediate directory, then the leaf). Each `Setup` receives the same `*Request` and can modify it incrementally.

## Run Resolution

`func Run` is defined **only at the root** `SETUP.md`. All leaves share the same `Run` function. Descendants must not redefine it.

If two test scenarios need a different `Run`, they require separate test trees, each rooted at its own `DOCTEST.md`.

## Type and Helper Scoping

- `type Request`, `type Response`, and `func Run` are defined once at the root. Child `SETUP.md` files **must not** redefine them.
- Helper functions in ancestor SETUPs are available to descendants. A child **must not** redefine a helper with the same name.

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