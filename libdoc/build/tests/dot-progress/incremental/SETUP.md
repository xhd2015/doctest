## Steps
1. No additional setup needed — the root Run handles everything.

```go
func Setup(t *testing.T, req *Request) error {
	_ = req
	return nil
}
```
