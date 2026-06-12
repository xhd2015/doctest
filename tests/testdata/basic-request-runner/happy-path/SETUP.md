# Happy Path Setup

## Steps

1. Select the default greeting action.
2. Override the request name for this leaf.

```go
func Setup(t *testing.T, req *Request) error {
	req.Action = "greet"
	req.Name = "runner"
	return nil
}
```
