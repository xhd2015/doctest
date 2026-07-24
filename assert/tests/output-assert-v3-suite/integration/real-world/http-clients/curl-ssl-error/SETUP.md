# Scenario

**Feature**: curl ssl error

```
# curl
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"curl: \\(60\\) SSL certificate problem: unable to get local issuer certificate",
	)
	req.Actual = "curl: (60) SSL certificate problem: unable to get local issuer certificate"
	return nil
}
```
