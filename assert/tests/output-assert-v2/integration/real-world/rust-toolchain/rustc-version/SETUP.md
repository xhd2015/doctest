# Scenario

**Feature**: rustc --version

```
# rustc
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__VER__: 'type=string, example=1.78.0'\n",
		"rustc 1.78.0",
	)
	req.Actual = "rustc 1.78.0"
	return nil
}
```
