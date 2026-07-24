# Scenario

**Feature**: go get

```
# go get
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__MOD__: 'type=string, example=github.com/x/y'\n",
		"go: added github\\.com/x/y v1\\.0\\.0",
	)
	req.Actual = "go: added github.com/x/y v1.0.0"
	return nil
}
```
