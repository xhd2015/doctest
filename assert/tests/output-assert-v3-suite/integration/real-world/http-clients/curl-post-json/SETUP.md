# Scenario

**Feature**: curl -d

```
# curl -X POST
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"\\{\"id\":1\\}",
	)
	req.Actual = "{\"id\":1}"
	return nil
}
```
