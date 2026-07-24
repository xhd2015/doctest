# Scenario

**Feature**: swift build

```
# swift build
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"Build complete! \\(1\\.23s\\)",
	)
	req.Actual = "Build complete! (1.23s)"
	return nil
}
```
