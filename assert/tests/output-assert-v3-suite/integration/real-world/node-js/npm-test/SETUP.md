# Scenario

**Feature**: npm test

```
# npm test
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"> test\nPASS",
	)
	req.Actual = "> test\nPASS"
	return nil
}
```
