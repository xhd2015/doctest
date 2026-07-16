# Scenario

**Feature**: bash -c echo

```
# bash
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"hi",
	)
	req.Actual = "hi"
	return nil
}
```
