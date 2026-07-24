# Scenario

**Feature**: git merge conflict

```
# git merge
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"CONFLICT \\(content\\): Merge conflict in main\\.go",
	)
	req.Actual = "CONFLICT (content): Merge conflict in main.go"
	return nil
}
```
