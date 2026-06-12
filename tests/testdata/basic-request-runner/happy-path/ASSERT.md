# Happy Path Assert

## Expected

- The inherited root `Run` returns a greeting response.
- No run error is returned.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected run error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected response")
	}
	if resp.Message != "hello runner" {
		t.Fatalf("message = %q, want %q", resp.Message, "hello runner")
	}
}
```
