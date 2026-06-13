## Expected
- `resp.BoolResult` is `true` (pattern ends with `/...` and does not start with `/`).

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if !resp.BoolResult {
		t.Fatalf("expected true for ./foo/..., got false")
	}
}
```
