# Scenario

**Feature**: rustc --version

```
# rustc
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__VER__: 'type=string, example=1.78.0'\n",
		"rustc 1\\.78\\.0",
	)
	req.Actual = "rustc 1.78.0"
	return nil
}
```
