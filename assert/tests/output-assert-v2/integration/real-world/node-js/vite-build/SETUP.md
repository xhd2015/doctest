# Scenario

**Feature**: vite build

```
# vite build
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"✓ built in 500ms.",
	)
	req.Actual = "✓ built in 500ms."
	return nil
}
```
