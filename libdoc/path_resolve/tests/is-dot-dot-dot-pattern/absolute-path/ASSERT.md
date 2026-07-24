## Expected
- `resp.BoolResult` is `true` for an absolute path ending in the dot-dot-dot suffix.

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if !resp.BoolResult {
		t.Fatalf("expected true for /foo/..., got false")
	}
}
```