# Expected Error Setup

## Steps

1. Select an action that the root `Run` intentionally reports as a run error.

```go
func Setup(t *testing.T, req *Request) error {
	req.Action = "fail"
	req.Name = "case"
	return nil
}
```
