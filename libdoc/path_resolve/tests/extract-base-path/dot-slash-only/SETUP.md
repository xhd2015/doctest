## Steps
- Input is `"./..."` (root directory pattern).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Input = "./..."
	return nil
}
```
