# Scenario

**Feature**: sbt test

```
# sbt test
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"\\[info\\] All tests passed\\.",
	)
	req.Actual = "[info] All tests passed."
	return nil
}
```
