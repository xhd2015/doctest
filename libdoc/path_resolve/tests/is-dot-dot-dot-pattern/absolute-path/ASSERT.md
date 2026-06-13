## Expected
- `resp.BoolResult` is `false` (starts with `/`, which is explicitly rejected).

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.BoolResult {
		t.Fatalf("expected false for /foo/..., got true")
	}
}
```
