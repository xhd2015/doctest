# Scenario

**Feature**: mypy

```
# mypy
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"main.py:1: error: Incompatible types",
	)
	req.Actual = "main.py:1: error: Incompatible types"
	return nil
}
```
