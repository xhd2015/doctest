# Scenario

**Feature**: ruff check

```
# ruff
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__N__: 'type=number, example=1'\n",
		"Found 1 error\\.",
	)
	req.Actual = "Found 1 error."
	return nil
}
```
