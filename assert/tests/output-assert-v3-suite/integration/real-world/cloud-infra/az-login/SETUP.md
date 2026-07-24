# Scenario

**Feature**: az login

```
# az login
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"You have logged in\\. Let us set up your subscription\\.",
	)
	req.Actual = "You have logged in. Let us set up your subscription."
	return nil
}
```
