## Steps
- Input is `"foo/bar"` (a plain path without `...`).

```go
func Setup(t *testing.T, req *Request) error {
	req.Input = "foo/bar"
	return nil
}
```
