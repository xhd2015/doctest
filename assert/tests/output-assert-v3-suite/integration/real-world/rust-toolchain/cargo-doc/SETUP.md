# Scenario

**Feature**: cargo doc

```
# cargo doc
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		" Documenting myapp v0\\.1\\.0",
	)
	req.Actual = " Documenting myapp v0.1.0"
	return nil
}
```
