## Expected
- `resp.BoolResult` is `false` (no `/...` suffix).

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.BoolResult {
		t.Fatalf("expected false for foo/bar, got true")
	}
}
```
