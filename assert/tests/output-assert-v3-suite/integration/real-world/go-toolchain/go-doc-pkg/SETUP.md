# Scenario

**Feature**: go doc

```
# go doc fmt
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"package fmt",
	)
	req.Actual = "package fmt"
	return nil
}
```
