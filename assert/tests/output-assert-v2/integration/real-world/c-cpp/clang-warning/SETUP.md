# Scenario

**Feature**: clang warning

```
# clang
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__LINE__: 'type=string, example=warning: unused variable ''x'' [-Wunused-variable]'\n",
		"__LINE__",
	)
	req.Actual = "warning: unused variable 'x' [-Wunused-variable]"
	return nil
}
```
