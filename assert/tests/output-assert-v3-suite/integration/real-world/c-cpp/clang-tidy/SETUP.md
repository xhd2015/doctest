# Scenario

**Feature**: clang-tidy

```
# clang-tidy
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"warning: use auto when initializing \\[modernize-use-auto\\]",
	)
	req.Actual = "warning: use auto when initializing [modernize-use-auto]"
	return nil
}
```
