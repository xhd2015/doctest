# Scenario

**Feature**: docker logs

```
# docker logs
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"listening on :8080",
	)
	req.Actual = "listening on :8080"
	return nil
}
```
