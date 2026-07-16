# Scenario

**Feature**: git stash

```
# git stash
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"Saved working directory and index state WIP on main",
	)
	req.Actual = "Saved working directory and index state WIP on main"
	return nil
}
```
