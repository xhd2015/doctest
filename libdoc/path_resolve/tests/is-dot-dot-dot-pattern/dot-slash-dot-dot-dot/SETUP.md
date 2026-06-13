## Steps
- Input is `"./foo/..."`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Input = "./foo/..."
	return nil
}
```
