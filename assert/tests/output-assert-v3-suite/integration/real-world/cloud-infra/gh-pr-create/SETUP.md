# Scenario

**Feature**: gh pr create

```
# gh pr create
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__URL__: 'type=string, example=https://github.com/x/y/pull/1'\n",
		"https://github\\.com/x/y/pull/1",
	)
	req.Actual = "https://github.com/x/y/pull/1"
	return nil
}
```
