## Steps
- Input is `"foo/..."` (without `./` prefix).

```go
func Setup(t *testing.T, req *Request) error {
	req.Input = "foo/..."
	return nil
}
```
