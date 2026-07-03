# Scenario

**Feature**: go work sync

```
# go work sync
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"go: syncing workspace",
	)
	req.Actual = "go: syncing workspace"
	return nil
}
```
