## Expected
- `resp.StringResult` is `"foo"` (strips `/...` and leading `./`).

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.StringResult != "foo" {
		t.Fatalf("expected 'foo', got %q", resp.StringResult)
	}
}
```
