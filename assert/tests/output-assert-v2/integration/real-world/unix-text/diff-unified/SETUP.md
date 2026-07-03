# Scenario

**Feature**: diff -u

```
# diff
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"--- a\n+++ b",
	)
	req.Actual = "--- a\n+++ b"
	return nil
}
```
