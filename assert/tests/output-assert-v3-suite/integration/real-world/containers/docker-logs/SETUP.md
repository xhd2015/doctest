# Scenario

**Feature**: docker logs

```
# docker logs
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"listening on :8080",
	)
	req.Actual = "listening on :8080"
	return nil
}
```
