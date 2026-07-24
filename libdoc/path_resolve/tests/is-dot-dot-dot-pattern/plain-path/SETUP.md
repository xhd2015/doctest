## Steps
- Input is `"foo/bar"` (a plain path without `...`).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Input = "foo/bar"
	return nil
}
```
