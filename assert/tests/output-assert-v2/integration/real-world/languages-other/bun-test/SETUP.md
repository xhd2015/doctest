# Scenario

**Feature**: bun test

```
# bun test
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__N__: 'type=number, example=2'\n",
		" 2 pass\n 0 fail",
	)
	req.Actual = " 2 pass\n 0 fail"
	return nil
}
```
