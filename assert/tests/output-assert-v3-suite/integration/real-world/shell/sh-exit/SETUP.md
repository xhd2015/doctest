# Scenario

**Feature**: sh exit code

```
# sh
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__CODE__: 'type=number, example=1'\n",
		"exit 1",
	)
	req.Actual = "exit 1"
	return nil
}
```
