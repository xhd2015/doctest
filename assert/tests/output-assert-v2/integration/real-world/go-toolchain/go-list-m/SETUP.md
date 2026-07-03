# Scenario

**Feature**: go list -m

```
# go list -m all
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__MOD__: 'type=string, example=example.com/x'\n",
		"example.com/x",
	)
	req.Actual = "example.com/x"
	return nil
}
```
