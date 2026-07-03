# Scenario

**Feature**: go mod init

```
# go mod init
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__MOD__: 'type=string, example=example.com/x'\n",
		"go: creating new go.mod: module example.com/x",
	)
	req.Actual = "go: creating new go.mod: module example.com/x"
	return nil
}
```
