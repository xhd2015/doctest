# Scenario

**Feature**: diff -u

```
# diff
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"--- a\n\\+\\+\\+ b",
	)
	req.Actual = "--- a\n+++ b"
	return nil
}
```
