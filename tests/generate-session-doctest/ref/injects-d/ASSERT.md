## Expected

- Leaf source satisfies inject contract (session.Doctest, d construct, pass d, no Chdir, no free vars).
- Root package source does not declare/assign free `DOCTEST_ROOT` / `DOCTEST_SESSION_ID`.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("ref assemble failed: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	assertInjectContract(t, "ref-leaf", resp.Source)
	if resp.RootSrc == "" {
		t.Fatal("expected non-empty ref root package source")
	}
	if hasPackageFreeDoctestVars(resp.RootSrc) {
		t.Fatalf("ref root package must not have free DOCTEST_* vars\n%s", resp.RootSrc)
	}
}
```
