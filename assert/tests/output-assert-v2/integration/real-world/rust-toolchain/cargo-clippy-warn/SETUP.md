# Scenario

**Feature**: cargo clippy

```
# cargo clippy
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"warning: unused variable: `x`",
	)
	req.Actual = "warning: unused variable: `x`"
	return nil
}
```
