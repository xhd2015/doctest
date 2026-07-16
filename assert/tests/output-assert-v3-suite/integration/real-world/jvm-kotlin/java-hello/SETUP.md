# Scenario

**Feature**: java Hello

```
# java
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"Hello",
	)
	req.Actual = "Hello"
	return nil
}
```
