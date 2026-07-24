# Scenario

**Feature**: podman ps

```
# podman ps
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"CONTAINER ID  IMAGE       COMMAND",
	)
	req.Actual = "CONTAINER ID  IMAGE       COMMAND"
	return nil
}
```
