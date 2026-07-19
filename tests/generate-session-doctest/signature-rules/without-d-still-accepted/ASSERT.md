## Expected

- Legacy without-d signatures still parse and pass rules (`ParseErr` empty).

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected run error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.ParseErr != "" {
		t.Fatalf("expected without-d signatures still accepted, got: %s", resp.ParseErr)
	}
}
```
