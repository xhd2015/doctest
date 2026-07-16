# Scenario

**Feature**: gradle test

```
# gradle test
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__N__: 'type=number, example=5'\n",
		"5 tests completed",
	)
	req.Actual = "5 tests completed"
	return nil
}
```
