# Scenario

**Feature**: ninja

```
# ninja
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__N__: 'type=number, example=10'\n",
		"\\[10/10\\] Linking CXX executable app",
	)
	req.Actual = "[10/10] Linking CXX executable app"
	return nil
}
```
