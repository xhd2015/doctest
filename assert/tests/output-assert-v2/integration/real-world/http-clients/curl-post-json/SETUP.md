# Scenario

**Feature**: curl -d

```
# curl -X POST
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"{\"id\":1}",
	)
	req.Actual = "{\"id\":1}"
	return nil
}
```
