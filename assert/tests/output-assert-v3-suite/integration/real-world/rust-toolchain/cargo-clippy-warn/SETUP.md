# Scenario

**Feature**: cargo clippy

```
# cargo clippy
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"warning: unused variable: `x`",
	)
	req.Actual = "warning: unused variable: `x`"
	return nil
}
```
