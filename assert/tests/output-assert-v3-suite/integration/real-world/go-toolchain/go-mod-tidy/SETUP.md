# Scenario

**Feature**: go mod tidy

```
# go mod tidy
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"go: downloading example\\.com/lib v1\\.0\\.0",
	)
	req.Actual = "go: downloading example.com/lib v1.0.0"
	return nil
}
```
