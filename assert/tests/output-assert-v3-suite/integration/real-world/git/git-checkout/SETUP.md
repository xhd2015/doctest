# Scenario

**Feature**: git checkout

```
# git checkout
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__BR__: 'type=string, example=dev'\n",
		"Switched to branch '__BR__",
	)
	req.Actual = "Switched to branch 'dev'"
	return nil
}
```
