# Scenario

**Feature**: git commit

```
# git commit
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__LINE__: 'type=string, example=[main abc1234] msg'\n",
		"__LINE__",
	)
	req.Actual = "[main abc1234] msg"
	return nil
}
```
