# Scenario

**Feature**: go mod init

```
# go mod init
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__MOD__: 'type=string, example=example.com/x'\n",
		"go: creating new go\\.mod: module example\\.com/x",
	)
	req.Actual = "go: creating new go.mod: module example.com/x"
	return nil
}
```
