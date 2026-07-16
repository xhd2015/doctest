# Scenario

**Feature**: uv sync

```
# uv sync
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"Resolved 10 packages in 12ms",
	)
	req.Actual = "Resolved 10 packages in 12ms"
	return nil
}
```
