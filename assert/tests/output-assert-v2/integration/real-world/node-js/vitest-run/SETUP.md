# Scenario

**Feature**: vitest run

```
# vitest
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		" Test Files  1 passed (1)",
	)
	req.Actual = " Test Files  1 passed (1)"
	return nil
}
```
