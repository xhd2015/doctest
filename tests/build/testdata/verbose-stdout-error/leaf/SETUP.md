## Steps

- Set the request name.

```go
func Setup(t *testing.T, req *Request) error {
    req.Name = "test"
    return nil
}
```
