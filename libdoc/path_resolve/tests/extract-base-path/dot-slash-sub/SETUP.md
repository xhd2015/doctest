## Steps
- Input is `"./foo/..."` (subdirectory pattern with dot-slash prefix).

```go
func Setup(t *testing.T, req *Request) error {
	req.Input = "./foo/..."
	return nil
}
```
