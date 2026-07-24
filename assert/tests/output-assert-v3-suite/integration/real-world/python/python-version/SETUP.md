# Scenario

**Feature**: python --version

```
# python
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__VER__: 'type=string, example=3.12.0'\n",
		"Python 3\\.12\\.0",
	)
	req.Actual = "Python 3.12.0"
	return nil
}
```
