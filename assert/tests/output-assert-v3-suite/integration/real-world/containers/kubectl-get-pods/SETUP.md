# Scenario

**Feature**: kubectl get pods

```
# kubectl get
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__NAME__: 'type=string, example=app-1'\n",
		"NAME   READY\napp-1   1/1",
	)
	req.Actual = "NAME   READY\napp-1   1/1"
	return nil
}
```
