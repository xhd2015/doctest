## Expected

- Generated source includes `_ *session.Doctest` as the second parameter on
  Setup/Run/Assert closures (author omitted the name).
- Call sites still pass the constructed `d`.
- Inject contract (no Chdir / free vars) still holds.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("classic assemble failed: %v", err)
	}
	assertInjectContract(t, "classic-omit-d", resp.Source)
	if !hasUnderscoreDoctestParam(resp.Source) {
		t.Fatalf("expected `_ *session.Doctest` when author omits d param\n%s", resp.Source)
	}
}
```
