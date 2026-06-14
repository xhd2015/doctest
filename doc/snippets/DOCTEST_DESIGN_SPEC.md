## Code blocks

Both `SETUP.md` and `ASSERT.md` may contain ```go...``` go code blocks. 

## DOCTEST.md

`DOCTEST.md` marks the root of the test tree. The whole test tree rooted from where `DOCTEST.md` begins forms a large decision tree. The root `SETUP.md` must define `type Request` and `type Response` — these types are shared by all descendants.

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