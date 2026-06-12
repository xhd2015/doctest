# Override Run Assert

## Expected

- The leaf `Run` is selected instead of the root `Run`.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected run error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected response")
	}
	if resp.Message != "override run for leaf" {
		t.Fatalf("message = %q", resp.Message)
	}
}
```
