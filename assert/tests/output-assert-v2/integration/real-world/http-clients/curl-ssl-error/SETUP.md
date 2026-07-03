# Scenario

**Feature**: curl ssl error

```
# curl
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"curl: (60) SSL certificate problem: unable to get local issuer certificate",
	)
	req.Actual = "curl: (60) SSL certificate problem: unable to get local issuer certificate"
	return nil
}
```
