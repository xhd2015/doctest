## Expected
- `resp.StringResult` is `"/foo/bar"` (strips `/...`, keeps absolute base path).

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.StringResult != "/foo/bar" {
		t.Fatalf("expected '/foo/bar', got %q", resp.StringResult)
	}
}
```