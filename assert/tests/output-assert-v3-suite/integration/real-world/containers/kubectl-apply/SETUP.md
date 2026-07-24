# Scenario

**Feature**: kubectl apply

```
# kubectl apply
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"deployment\\.apps/app created",
	)
	req.Actual = "deployment.apps/app created"
	return nil
}
```
