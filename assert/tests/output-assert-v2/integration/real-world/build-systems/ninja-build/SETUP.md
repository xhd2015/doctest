# Scenario

**Feature**: ninja

```
# ninja
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__LINE__: 'type=string, example=[10/10] Linking CXX executable app'\n",
		"__LINE__",
	)
	req.Actual = "[10/10] Linking CXX executable app"
	return nil
}
```
