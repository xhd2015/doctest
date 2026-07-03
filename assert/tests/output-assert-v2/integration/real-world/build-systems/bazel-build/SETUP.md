# Scenario

**Feature**: bazel build

```
# bazel build
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"INFO: Build completed successfully, 3 total actions",
	)
	req.Actual = "INFO: Build completed successfully, 3 total actions"
	return nil
}
```
