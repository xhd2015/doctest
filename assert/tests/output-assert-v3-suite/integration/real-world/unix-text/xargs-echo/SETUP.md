# Scenario

**Feature**: xargs

```
# xargs echo
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"one two",
	)
	req.Actual = "one two"
	return nil
}
```
