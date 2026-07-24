# Scenario

**Feature**: vite build

```
# vite build
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"✓ built in 500ms\\.",
	)
	req.Actual = "✓ built in 500ms."
	return nil
}
```
