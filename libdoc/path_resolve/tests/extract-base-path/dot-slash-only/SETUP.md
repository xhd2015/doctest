## Steps
- Input is `"./..."` (root directory pattern).

```go
func Setup(t *testing.T, req *Request) error {
	req.Input = "./..."
	return nil
}
```
