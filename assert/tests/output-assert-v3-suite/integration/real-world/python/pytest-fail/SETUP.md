# Scenario

**Feature**: pytest fail

```
# pytest
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"FAILED test_a\\.py::test_x - AssertionError",
	)
	req.Actual = "FAILED test_a.py::test_x - AssertionError"
	return nil
}
```
