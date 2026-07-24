# Scenario

**Feature**: openssl s_client

```
# openssl
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"CONNECTED\\(00000003\\)",
	)
	req.Actual = "CONNECTED(00000003)"
	return nil
}
```
