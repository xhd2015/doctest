# Scenario

**Feature**: go get

```
# go get
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__MOD__: 'type=string, example=github.com/x/y'\n",
		"go: added github.com/x/y v1.0.0",
	)
	req.Actual = "go: added github.com/x/y v1.0.0"
	return nil
}
```
