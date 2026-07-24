# Scenario

**Feature**: glab mr list

```
# glab mr list
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__IID__: 'type=number, example=1'\n",
		"!1 Draft: feature",
	)
	req.Actual = "!1 Draft: feature"
	return nil
}
```
