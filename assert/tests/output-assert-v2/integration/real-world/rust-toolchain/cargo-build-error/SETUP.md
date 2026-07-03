# Scenario

**Feature**: cargo build error

```
# cargo build
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__LINE__: 'type=string, example=error[E0425]: cannot find value `X` in this scope'\n",
		"__LINE__",
	)
	req.Actual = "error[E0425]: cannot find value `X` in this scope"
	return nil
}
```
