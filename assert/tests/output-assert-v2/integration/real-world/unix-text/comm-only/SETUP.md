# Scenario

**Feature**: comm -13

```
# comm
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"only-b",
	)
	req.Actual = "only-b"
	return nil
}
```
