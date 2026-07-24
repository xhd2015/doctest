## Steps
- Input is `"..."` (bare dots, no leading `./`).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Input = "..."
	return nil
}
```
