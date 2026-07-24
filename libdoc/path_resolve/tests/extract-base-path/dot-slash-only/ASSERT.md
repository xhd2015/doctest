## Expected
- `resp.StringResult` is `"."` (strips `/...` and `./`, left with empty → returns `"."`).

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.StringResult != "." {
		t.Fatalf("expected '.', got %q", resp.StringResult)
	}
}
```
