# Scenario

**Feature**: jest

```
# jest
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__N__: 'type=number, example=5'\n",
		"Tests:       5 passed, 5 total",
	)
	req.Actual = "Tests:       5 passed, 5 total"
	return nil
}
```
