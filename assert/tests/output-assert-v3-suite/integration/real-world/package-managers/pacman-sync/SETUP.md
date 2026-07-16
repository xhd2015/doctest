# Scenario

**Feature**: pacman -Sy

```
# pacman
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		":: Synchronizing package databases\\.\\.\\.",
	)
	req.Actual = ":: Synchronizing package databases..."
	return nil
}
```
