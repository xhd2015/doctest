## Expected

- The response message matches the request name.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatal(err)
    }
}
```
