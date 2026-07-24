# Scenario

**Feature**: mypy

```
# mypy
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"main\\.py:1: error: Incompatible types",
	)
	req.Actual = "main.py:1: error: Incompatible types"
	return nil
}
```
