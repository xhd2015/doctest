## Expected

- The response is correct.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatal(err)
    }
}
```
