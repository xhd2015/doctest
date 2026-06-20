## Expected

- `resp.Display` is `"."`.

## Errors

- `err` is nil.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.Display != "." {
		t.Fatalf("expected \".\", got %q (cwd=%q path=%q)", resp.Display, resp.Cwd, req.Path)
	}
}
```