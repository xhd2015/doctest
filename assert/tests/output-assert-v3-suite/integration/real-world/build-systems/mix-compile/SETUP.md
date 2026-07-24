# Scenario

**Feature**: mix compile

```
# mix compile
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"Compiling 1 file \\(\\.ex\\)",
	)
	req.Actual = "Compiling 1 file (.ex)"
	return nil
}
```
