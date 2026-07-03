# Scenario

**Feature**: dash test

```
# dash
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"ok",
	)
	req.Actual = "ok"
	return nil
}
```
