# Scenario

**Feature**: sbt compile

```
# sbt compile
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__LINE__: 'type=string, example=[success] Total time: 5 s, completed'\n",
		"__LINE__",
	)
	req.Actual = "[success] Total time: 5 s, completed"
	return nil
}
```
