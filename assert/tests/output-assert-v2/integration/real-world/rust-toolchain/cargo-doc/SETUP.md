# Scenario

**Feature**: cargo doc

```
# cargo doc
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		" Documenting myapp v0.1.0",
	)
	req.Actual = " Documenting myapp v0.1.0"
	return nil
}
```
