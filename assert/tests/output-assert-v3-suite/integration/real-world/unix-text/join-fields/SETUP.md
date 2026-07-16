# Scenario

**Feature**: join

```
# join
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"1 a",
	)
	req.Actual = "1 a"
	return nil
}
```
