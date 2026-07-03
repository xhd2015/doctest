# Scenario

**Feature**: git merge conflict

```
# git merge
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"CONFLICT (content): Merge conflict in main.go",
	)
	req.Actual = "CONFLICT (content): Merge conflict in main.go"
	return nil
}
```
