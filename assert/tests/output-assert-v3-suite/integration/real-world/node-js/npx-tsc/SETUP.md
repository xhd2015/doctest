# Scenario

**Feature**: npx tsc

```
# npx tsc
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"Version 5\\.0\\.0",
	)
	req.Actual = "Version 5.0.0"
	return nil
}
```
