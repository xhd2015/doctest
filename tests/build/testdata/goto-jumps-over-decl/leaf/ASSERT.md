## Expected

- Inner leaf would pass if the generated suite compiled.

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
```
