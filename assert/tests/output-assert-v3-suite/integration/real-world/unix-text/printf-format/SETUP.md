# Scenario

**Feature**: printf

```
# printf
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"value=1",
	)
	req.Actual = "value=1"
	return nil
}
```
