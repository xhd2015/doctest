## Expected
- `resp.ErrResult` is non-empty (error because no dirs found).

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err == nil {
		t.Fatal("expected error for no dirs")
	}
	if resp.ErrResult == "" {
		t.Fatal("expected ErrResult to contain error message")
	}
}
```
