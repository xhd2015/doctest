# Scenario

**Feature**: git log

```
# git log
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__HASH__: 'type=string, example=abc1234'\n",
		"abc1234 init",
	)
	req.Actual = "abc1234 init"
	return nil
}
```
