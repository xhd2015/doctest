# Scenario

**Feature**: kubectl describe

```
# kubectl describe
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__NS__: 'type=string, example=default'\n",
		"Namespace:        default",
	)
	req.Actual = "Namespace:        default"
	return nil
}
```
