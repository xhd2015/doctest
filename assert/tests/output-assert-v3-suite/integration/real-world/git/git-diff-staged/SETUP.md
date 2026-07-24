# Scenario

**Feature**: git diff --cached

```
# git diff
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"diff --git a/a b/a",
	)
	req.Actual = "diff --git a/a b/a"
	return nil
}
```
