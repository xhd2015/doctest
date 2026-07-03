# Scenario

**Feature**: go env GOMOD

```
# go env
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"/tmp/go.mod",
	)
	req.Actual = "/tmp/go.mod"
	return nil
}
```
