# Scenario

**Feature**: scalac

```
# scalac
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"compiling 1 source file to target",
	)
	req.Actual = "compiling 1 source file to target"
	return nil
}
```
