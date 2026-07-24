# Scenario

**Feature**: go test fail

```
# go test
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"--- FAIL: TestX \\(0\\.00s\\)",
	)
	req.Actual = "--- FAIL: TestX (0.00s)"
	return nil
}
```
