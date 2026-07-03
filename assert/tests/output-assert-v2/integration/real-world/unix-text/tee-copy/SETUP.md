# Scenario

**Feature**: tee

```
# tee
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"copied line",
	)
	req.Actual = "copied line"
	return nil
}
```
