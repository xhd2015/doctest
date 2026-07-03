# Scenario

**Feature**: glab mr list

```
# glab mr list
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__IID__: 'type=number, example=1'\n",
		"!1 Draft: feature",
	)
	req.Actual = "!1 Draft: feature"
	return nil
}
```
