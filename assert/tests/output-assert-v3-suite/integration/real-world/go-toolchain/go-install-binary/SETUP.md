# Scenario

**Feature**: go install

```
# go install
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__BIN__: 'type=string, example=mytool'\n",
		"go: installing mytool",
	)
	req.Actual = "go: installing mytool"
	return nil
}
```
