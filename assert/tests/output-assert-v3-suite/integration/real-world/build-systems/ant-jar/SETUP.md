# Scenario

**Feature**: ant jar

```
# ant jar
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"BUILD SUCCESSFUL\nTotal time: 2 seconds",
	)
	req.Actual = "BUILD SUCCESSFUL\nTotal time: 2 seconds"
	return nil
}
```
