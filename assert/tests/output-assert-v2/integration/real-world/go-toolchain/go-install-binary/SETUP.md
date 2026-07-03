# Scenario

**Feature**: go install

```
# go install
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__BIN__: 'type=string, example=mytool'\n",
		"go: installing mytool",
	)
	req.Actual = "go: installing mytool"
	return nil
}
```
