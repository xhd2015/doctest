# Scenario

**Feature**: redis-cli PING

```
# redis-cli
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"PONG",
	)
	req.Actual = "PONG"
	return nil
}
```
