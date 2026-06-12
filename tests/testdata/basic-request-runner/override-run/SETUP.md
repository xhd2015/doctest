# Override Run Setup

## Steps

1. Mutate the inherited request.
2. Define a deeper `Run` for this leaf.

```go
func Setup(t *testing.T, req *Request) error {
	req.Action = "ignored-by-override"
	req.Name = "leaf"
	return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
	return &Response{Message: "override run for " + req.Name}, nil
}
```
