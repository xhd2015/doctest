# Scenario

**Feature**: java Hello

```
# java
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"Hello",
	)
	req.Actual = "Hello"
	return nil
}
```
