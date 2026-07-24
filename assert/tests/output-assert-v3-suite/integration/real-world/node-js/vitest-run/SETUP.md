# Scenario

**Feature**: vitest run

```
# vitest
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		" Test Files  1 passed \\(1\\)",
	)
	req.Actual = " Test Files  1 passed (1)"
	return nil
}
```
