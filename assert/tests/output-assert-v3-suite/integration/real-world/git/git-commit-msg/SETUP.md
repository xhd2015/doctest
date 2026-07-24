# Scenario

**Feature**: git commit

```
# git commit
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__HASH__: 'type=string, example=abc1234'\n",
		"\\[main abc1234\\] msg",
	)
	req.Actual = "[main abc1234] msg"
	return nil
}
```
