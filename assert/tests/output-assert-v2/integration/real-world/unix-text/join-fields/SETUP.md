# Scenario

**Feature**: join

```
# join
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"1 a",
	)
	req.Actual = "1 a"
	return nil
}
```
