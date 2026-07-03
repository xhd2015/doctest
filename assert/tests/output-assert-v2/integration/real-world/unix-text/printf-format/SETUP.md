# Scenario

**Feature**: printf

```
# printf
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"value=1",
	)
	req.Actual = "value=1"
	return nil
}
```
