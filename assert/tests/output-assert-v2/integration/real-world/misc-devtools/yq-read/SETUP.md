# Scenario

**Feature**: yq

```
# yq
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__VAL__: 'type=string, example=prod'\n",
		"prod",
	)
	req.Actual = "prod"
	return nil
}
```
