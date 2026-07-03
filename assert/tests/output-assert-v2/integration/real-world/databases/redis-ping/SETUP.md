# Scenario

**Feature**: redis-cli PING

```
# redis-cli
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"PONG",
	)
	req.Actual = "PONG"
	return nil
}
```
