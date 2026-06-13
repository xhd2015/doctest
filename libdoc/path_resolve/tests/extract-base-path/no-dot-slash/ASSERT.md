## Expected
- `resp.StringResult` is `"foo"` (strips `/...`, no leading `./` to strip).

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.StringResult != "foo" {
		t.Fatalf("expected 'foo', got %q", resp.StringResult)
	}
}
```
