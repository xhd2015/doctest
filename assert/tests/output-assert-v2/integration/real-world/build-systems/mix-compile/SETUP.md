# Scenario

**Feature**: mix compile

```
# mix compile
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"Compiling 1 file (.ex)",
	)
	req.Actual = "Compiling 1 file (.ex)"
	return nil
}
```
