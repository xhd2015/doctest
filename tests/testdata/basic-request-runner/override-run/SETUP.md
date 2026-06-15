# Override Run Setup

## Steps

1. Mutate the inherited request to use the root Run.

```go
func Setup(t *testing.T, req *Request) error {
	req.Action = "greet"
	req.Name = "leaf"
	return nil
}
```
