# Scenario

**Feature**: terraform apply

```
# terraform apply
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"Apply complete! Resources: 1 added, 0 changed, 0 destroyed.",
	)
	req.Actual = "Apply complete! Resources: 1 added, 0 changed, 0 destroyed."
	return nil
}
```
