# Scenario

**Feature**: cargo tree

```
# cargo tree
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__CRATE__: 'type=string, example=serde'\n",
		"serde v1.0.0",
	)
	req.Actual = "serde v1.0.0"
	return nil
}
```
