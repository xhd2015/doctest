## Code blocks

Both `SETUP.md` and `ASSERT.md` may contain ```go...``` go code blocks. 

Root `SETUP.md` defines `type Request` and `type Response` — shared by all descendants.

| Function | Where | Signature | Notes |
|----------|-------|-----------|-------|
| `Setup` | any SETUP.md | `(t *testing.T, req *Request) error` | called root→leaf before Run; body must not be stub |
| `Run` | any SETUP.md | `(t *testing.T, req *Request) (*Response, error)` | **deepest wins**; root provides stub so tests fail RED |
| `Assert` | every ASSERT.md | `(t *testing.T, req *Request, resp *Response, err error)` | fail via `t.Fatal`/`t.Fatalf` |

Import target package directly. For unexported functions, use **`TestExported_`** prefix:
`func TestExported_foo() { foo() }` — then `import "mypkg"; mypkg.TestExported_foo()` in the code block.

- Every `SETUP.md` must have a Go block as **final content**; child must not redefine Request/Response
- Signatures must match exactly; `func Setup` body must not be a stub (`return nil`)
- At least one `Run` in the chain; every `ASSERT.md` must have `func Assert`

## Test Fixture Data

Abstract fixture data into standalone files, not inline code.

- Single file → place alongside `ASSERT.md`
- Multiple files → place in `testdata/` alongside `ASSERT.md`

Code reads them with directly filename reference as each `ASSERT.md` runs in its own directory.

> Full spec, run: `doctest skill doc-spec show` && `doctest skill code-spec show`