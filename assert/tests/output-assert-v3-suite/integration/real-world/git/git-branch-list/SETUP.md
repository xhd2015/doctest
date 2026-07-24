# Scenario

**Feature**: git branch

```
# git branch
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__BR__: 'type=string, example=feature'\n",
		"\\* feature",
	)
	req.Actual = "* feature"
	return nil
}
```
