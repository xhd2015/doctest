# Scenario

**Feature**: tsc --noEmit

```
# tsc
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"Found 0 errors\\.",
	)
	req.Actual = "Found 0 errors."
	return nil
}
```
