# Scenario

**Feature**: crystal build

```
# crystal build
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"Compiling app",
	)
	req.Actual = "Compiling app"
	return nil
}
```
