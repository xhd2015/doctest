## Expected
- `resp.BoolResult` is `false` (bare `...` does not match the `/...` suffix requirement).

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.BoolResult {
		t.Fatalf("expected false for ..., got true")
	}
}
```
