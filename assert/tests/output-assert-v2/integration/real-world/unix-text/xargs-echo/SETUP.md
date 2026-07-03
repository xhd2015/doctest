# Scenario

**Feature**: xargs

```
# xargs echo
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"one two",
	)
	req.Actual = "one two"
	return nil
}
```
