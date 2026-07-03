# Scenario

**Feature**: eslint

```
# eslint
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__N__: 'type=number, example=2'\n",
		"✖ 2 problems (1 error, 1 warning)",
	)
	req.Actual = "✖ 2 problems (1 error, 1 warning)"
	return nil
}
```
