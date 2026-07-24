# Scenario

**Feature**: fish --version

```
# fish
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"fish, version 3\\.6\\.0",
	)
	req.Actual = "fish, version 3.6.0"
	return nil
}
```
