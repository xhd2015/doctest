# Scenario

**Feature**: valgrind

```
# valgrind
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"ERROR SUMMARY: 0 errors from 0 contexts",
	)
	req.Actual = "ERROR SUMMARY: 0 errors from 0 contexts"
	return nil
}
```
