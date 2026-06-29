## Steps
- Input is `"/foo/bar/..."` (absolute path with `/...` suffix).

```go
func Setup(t *testing.T, req *Request) error {
	req.Input = "/foo/bar/..."
	return nil
}
```