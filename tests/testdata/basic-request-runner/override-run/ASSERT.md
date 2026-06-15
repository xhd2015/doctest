# Override Run Assert

## Expected

- The root `Run` is executed (child cannot redefine Run).

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected run error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected response")
	}
	if resp.Message != "hello leaf" {
		t.Fatalf("message = %q", resp.Message)
	}
}
```
