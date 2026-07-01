## Expected
- Match fails. After P4 strips the leading empty line, the pattern is a pure
  single literal `foo` with no trailing newline; the actual ends with `\n`,
  so the strict trailing-newline policy is violated.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireMatchError(t, resp, "trailing newline")
}
```
