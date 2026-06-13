## Expected
- `len(resp.DirsResult)` is 1.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.DirsResult) != 1 {
		t.Fatalf("expected 1 dir, got %d: %v", len(resp.DirsResult), resp.DirsResult)
	}
}
```
