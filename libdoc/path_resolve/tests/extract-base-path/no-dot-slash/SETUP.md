## Steps
- Input is `"foo/..."` (without `./` prefix).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Input = "foo/..."
	return nil
}
```
