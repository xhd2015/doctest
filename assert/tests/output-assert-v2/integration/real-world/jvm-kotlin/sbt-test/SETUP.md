# Scenario

**Feature**: sbt test

```
# sbt test
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__LINE__: 'type=string, example=[info] All tests passed.'\n",
		"__LINE__",
	)
	req.Actual = "[info] All tests passed."
	return nil
}
```
