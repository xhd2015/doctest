# Scenario

**Feature**: git diff --cached

```
# git diff
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"diff --git a/a b/a",
	)
	req.Actual = "diff --git a/a b/a"
	return nil
}
```
