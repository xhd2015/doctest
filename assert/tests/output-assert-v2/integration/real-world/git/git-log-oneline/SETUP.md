# Scenario

**Feature**: git log

```
# git log
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__HASH__: 'type=string, example=abc1234'\n",
		"abc1234 init",
	)
	req.Actual = "abc1234 init"
	return nil
}
```
