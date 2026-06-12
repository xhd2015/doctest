## Steps

- Set the request name for this leaf case.

```go
func Setup(t *testing.T, req *Request) error {
    req.Name = "test"
    return nil
}
```
