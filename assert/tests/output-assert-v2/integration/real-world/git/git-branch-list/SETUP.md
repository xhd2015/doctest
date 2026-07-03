# Scenario

**Feature**: git branch

```
# git branch
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__BR__: 'type=string, example=feature'\n",
		"* feature",
	)
	req.Actual = "* feature"
	return nil
}
```
