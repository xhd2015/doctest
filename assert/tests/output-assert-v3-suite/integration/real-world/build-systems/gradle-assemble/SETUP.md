# Scenario

**Feature**: gradle assemble

```
# gradle
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"BUILD SUCCESSFUL in 3s",
	)
	req.Actual = "BUILD SUCCESSFUL in 3s"
	return nil
}
```
