# Scenario

**Feature**: node --version

```
# node -v
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__VER__: 'type=string, example=v20.0.0'\n",
		"v20.0.0",
	)
	req.Actual = "v20.0.0"
	return nil
}
```
