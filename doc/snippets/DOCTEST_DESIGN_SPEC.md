## Code blocks

Both `SETUP.md` and `ASSERT.md` may contain ```go...``` go code blocks. 

## DOCTEST.md

`DOCTEST.md` marks the root of the test tree. The whole test tree rooted from where `DOCTEST.md` begins forms a large decision tree. The root `SETUP.md` must define `type Request` and `type Response` — these types are shared by all descendants.

### Nested DOCTEST.md

A subdirectory that contains its own `DOCTEST.md` becomes a **self-contained test root**. The doctest runner stops walking at `DOCTEST.md` boundaries and treats each root independently — **no inheritance crosses a `DOCTEST.md` boundary**.

A nested root's `SETUP.md` must be entirely self-sufficient:
- It must define its own `Request`/`Response` types
- It must provide `Setup`/`Run` or let descendant SETUPs provide `Run`
- Any external binaries (e.g., the doctest binary for `req.Bin`) must be built or resolved within that root's own `Setup`
- The parent root's `Setup` is never executed for leaves under a nested DOCTEST.md
- Paths like `DOCTEST_ROOT/..` shift — from a deeper root, use `DOCTEST_ROOT/../..` to reach the module root

### Tree Organization

1. **Parent → child dirs**: scenarios become more concrete by narrowing one or a few params from `Request`.
2. **Sibling dirs**: must be mutually exclusive — each tests a different scenario branch.

## SETUP.md

Every `SETUP.md` must have a Go block as **final content**. Child must not redefine `Request`/`Response`.

| Function | Signature | Notes |
|----------|-----------|-------|
| `Setup` | `(t *testing.T, req *Request) error` | Called root→leaf before `Run`; body must not be stub |
| `Run` | `(t *testing.T, req *Request) (*Response, error)` | **Deepest wins**; root provides stub so tests fail RED |

At least one `Run` in the chain. Signatures must match exactly. `func Setup` body must not be a stub (`return nil`).

# Inheritance

## Setup Chain

`Setup` functions accumulate from root → leaf: every ancestor `SETUP.md`'s `func Setup` is called in order (root first, then each intermediate directory, then the leaf). Each `Setup` receives the same `*Request` and can modify it incrementally.

## Run Resolution

The `Run` function follows **deepest wins**: the last `func Run` defined in the chain (closest to the leaf) is the one that executes. A child `Run` completely replaces any ancestor `Run`.

## Type and Helper Scoping

- `type Request` and `type Response` are defined once at the root. Child `SETUP.md` files **must not** redefine them.
- Helper functions in ancestor SETUPs are available to descendants. A child **must not** redefine a helper with the same name.

## DOCTEST.md Boundary

A `DOCTEST.md` file creates an **inheritance firewall**. No code, types, helpers, or `Setup` functions cross a `DOCTEST.md` boundary. Each tree rooted at a `DOCTEST.md` is a self-contained decision tree with its own `Request`/`Response` types and its own setup chain.

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