# Scenario

**Feature**: swift build

```
# swift build
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"Build complete! (1.23s)",
	)
	req.Actual = "Build complete! (1.23s)"
	return nil
}
```
