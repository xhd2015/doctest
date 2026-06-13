## Steps
- Input is `"/foo/..."` (absolute path starting with `/`).

```go
func Setup(t *testing.T, req *Request) error {
	req.Input = "/foo/..."
	return nil
}
```
