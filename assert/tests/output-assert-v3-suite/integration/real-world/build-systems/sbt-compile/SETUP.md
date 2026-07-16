# Scenario

**Feature**: sbt compile

```
# sbt compile
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"\\[success\\] Total time: 5 s, completed",
	)
	req.Actual = "[success] Total time: 5 s, completed"
	return nil
}
```
