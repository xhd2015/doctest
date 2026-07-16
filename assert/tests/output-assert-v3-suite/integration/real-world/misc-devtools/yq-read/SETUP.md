# Scenario

**Feature**: yq

```
# yq
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__VAL__: 'type=string, example=prod'\n",
		"prod",
	)
	req.Actual = "prod"
	return nil
}
```
