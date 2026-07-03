# Scenario

**Feature**: valgrind

```
# valgrind
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"ERROR SUMMARY: 0 errors from 0 contexts",
	)
	req.Actual = "ERROR SUMMARY: 0 errors from 0 contexts"
	return nil
}
```
