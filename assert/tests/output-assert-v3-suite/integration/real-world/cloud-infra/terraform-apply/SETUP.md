# Scenario

**Feature**: terraform apply

```
# terraform apply
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"Apply complete! Resources: 1 added, 0 changed, 0 destroyed\\.",
	)
	req.Actual = "Apply complete! Resources: 1 added, 0 changed, 0 destroyed."
	return nil
}
```
