# Scenario

**Feature**: go env GOMOD

```
# go env
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"/tmp/go\\.mod",
	)
	req.Actual = "/tmp/go.mod"
	return nil
}
```
