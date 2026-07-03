# Scenario

**Feature**: tail -n

```
# tail file
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"line nine\nline ten",
	)
	req.Actual = "line nine\nline ten"
	return nil
}
```
