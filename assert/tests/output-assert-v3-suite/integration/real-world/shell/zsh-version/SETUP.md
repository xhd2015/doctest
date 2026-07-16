# Scenario

**Feature**: zsh --version

```
# zsh
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__VER__: 'type=string, example=5.9'\n",
		"zsh 5\\.9 \\(x86_64-apple-darwin23\\.0\\.0\\)",
	)
	req.Actual = "zsh 5.9 (x86_64-apple-darwin23.0.0)"
	return nil
}
```
