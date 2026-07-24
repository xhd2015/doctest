# Scenario

**Feature**: find -name

```
# find
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__PATH__: 'type=string, example=./a.go'\n",
		"\\./a\\.go",
	)
	req.Actual = "./a.go"
	return nil
}
```
