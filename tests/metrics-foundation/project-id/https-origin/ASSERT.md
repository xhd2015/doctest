## Expected

- `ProjectIDFromOrigin` returns exactly `github.com_xhd2015_doctest`.
- No scheme (`https`), no `.git` suffix, slashes replaced by `_`.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "github.com_xhd2015_doctest"
	if resp.ProjectID != want {
		t.Fatalf("project id = %q, want %q", resp.ProjectID, want)
	}
}
```
