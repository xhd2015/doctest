# Scenario

**Feature**: kubectl apply

```
# kubectl apply
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"deployment.apps/app created",
	)
	req.Actual = "deployment.apps/app created"
	return nil
}
```
