# Scenario

**Feature**: crystal build

```
# crystal build
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"Compiling app",
	)
	req.Actual = "Compiling app"
	return nil
}
```
