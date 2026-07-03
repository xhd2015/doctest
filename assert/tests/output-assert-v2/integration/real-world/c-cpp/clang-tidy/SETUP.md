# Scenario

**Feature**: clang-tidy

```
# clang-tidy
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__LINE__: 'type=string, example=warning: use auto when initializing [modernize-use-auto]'\n",
		"__LINE__",
	)
	req.Actual = "warning: use auto when initializing [modernize-use-auto]"
	return nil
}
```
