## Expected

- Generated classic source satisfies the full inject contract:
  - `session.Doctest` type and session import present
  - `d` constructed via `&session.Doctest{...}`
  - Setup / Run / Assert call sites pass `d`
  - no leaf Chdir / Getwd restore
  - no package free `DOCTEST_ROOT` / `DOCTEST_SESSION_ID`

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("classic assemble failed: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	assertInjectContract(t, "classic", resp.Source)
}
```
