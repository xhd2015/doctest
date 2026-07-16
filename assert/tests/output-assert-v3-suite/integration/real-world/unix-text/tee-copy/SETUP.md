# Scenario

**Feature**: tee

```
# tee
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"copied line",
	)
	req.Actual = "copied line"
	return nil
}
```
