# Scenario

**Feature**: go version

```
# go version
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__VER__: 'type=string, example=go1.22.0'\n",
		"go version go1.22.0 darwin/arm64",
	)
	req.Actual = "go version go1.22.0 darwin/arm64"
	return nil
}
```
