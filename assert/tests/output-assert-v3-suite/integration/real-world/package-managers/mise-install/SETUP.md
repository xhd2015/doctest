# Scenario

**Feature**: mise install

```
# mise install
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__TOOL__: 'type=string, example=node@20'\n",
		"mise node@20\\.0\\.0 installed",
	)
	req.Actual = "mise node@20.0.0 installed"
	return nil
}
```
