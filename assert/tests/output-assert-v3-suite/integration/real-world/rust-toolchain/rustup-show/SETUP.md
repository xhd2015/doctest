# Scenario

**Feature**: rustup show

```
# rustup
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__CH__: 'type=string, example=stable'\n",
		"active toolchain: stable-x86_64-unknown-linux-gnu \\(default\\)",
	)
	req.Actual = "active toolchain: stable-x86_64-unknown-linux-gnu (default)"
	return nil
}
```
