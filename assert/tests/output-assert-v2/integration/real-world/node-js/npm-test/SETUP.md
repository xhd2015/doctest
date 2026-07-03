# Scenario

**Feature**: npm test

```
# npm test
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"> test\nPASS",
	)
	req.Actual = "> test\nPASS"
	return nil
}
```
