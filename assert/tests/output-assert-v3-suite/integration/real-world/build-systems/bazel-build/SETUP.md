# Scenario

**Feature**: bazel build

```
# bazel build
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"INFO: Build completed successfully, 3 total actions",
	)
	req.Actual = "INFO: Build completed successfully, 3 total actions"
	return nil
}
```
