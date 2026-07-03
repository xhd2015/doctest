# Scenario

**Feature**: pulumi preview

```
# pulumi preview
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"Previewing update (dev)",
	)
	req.Actual = "Previewing update (dev)"
	return nil
}
```
