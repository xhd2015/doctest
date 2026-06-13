## Expected
- `len(resp.DirsResult)` is 3 (sub-a, sub-b, sub-c).

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.DirsResult) != 3 {
		t.Fatalf("expected 3 dirs, got %d: %v", len(resp.DirsResult), resp.DirsResult)
	}
}
```
