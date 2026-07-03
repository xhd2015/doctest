# Scenario

**Feature**: openssl s_client

```
# openssl
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"CONNECTED(00000003)",
	)
	req.Actual = "CONNECTED(00000003)"
	return nil
}
```
