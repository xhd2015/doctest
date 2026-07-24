## Steps
- Input is `"./foo/..."`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Input = "./foo/..."
	return nil
}
```
