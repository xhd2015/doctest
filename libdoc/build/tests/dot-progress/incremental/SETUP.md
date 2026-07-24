## Steps
1. No additional setup needed — the root Run handles everything.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = req
	return nil
}
```
