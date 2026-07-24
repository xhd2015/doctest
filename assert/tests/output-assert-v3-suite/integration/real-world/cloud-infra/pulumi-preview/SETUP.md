# Scenario

**Feature**: pulumi preview

```
# pulumi preview
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"Previewing update \\(dev\\)",
	)
	req.Actual = "Previewing update (dev)"
	return nil
}
```
