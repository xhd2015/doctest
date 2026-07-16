# Scenario

**Feature**: go work sync

```
# go work sync
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"go: syncing workspace",
	)
	req.Actual = "go: syncing workspace"
	return nil
}
```
