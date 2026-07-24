# Scenario

**Feature**: terraform plan

```
# terraform plan
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"Plan: 1 to add, 0 to change, 0 to destroy\\.",
	)
	req.Actual = "Plan: 1 to add, 0 to change, 0 to destroy."
	return nil
}
```
