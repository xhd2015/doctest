## Expected

- Generated source keeps author name `d *session.Doctest` on Setup/Run/Assert.
- Inject contract holds (construct, pass, no Chdir, no free vars).

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("classic assemble failed: %v", err)
	}
	assertInjectContract(t, "classic-named-d", resp.Source)
	if !hasNamedDoctestParam(resp.Source, "d") {
		t.Fatalf("expected kept name `d *session.Doctest`\n%s", resp.Source)
	}
}
```
