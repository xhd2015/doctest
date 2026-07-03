# Scenario

**Feature**: podman ps

```
# podman ps
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"CONTAINER ID  IMAGE       COMMAND",
	)
	req.Actual = "CONTAINER ID  IMAGE       COMMAND"
	return nil
}
```
