# Expected Error Assert

## Errors

- The run error is passed to `Assert`.
- Setup errors are not involved in this case.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err == nil {
		t.Fatal("expected run error")
	}
	if resp != nil {
		t.Fatalf("response = %#v, want nil", resp)
	}
	if err.Error() != "requested failure for case" {
		t.Fatalf("err = %q", err.Error())
	}
}
```
