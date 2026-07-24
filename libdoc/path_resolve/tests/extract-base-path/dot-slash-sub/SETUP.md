## Steps
- Input is `"./foo/..."` (subdirectory pattern with dot-slash prefix).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Input = "./foo/..."
	return nil
}
```
