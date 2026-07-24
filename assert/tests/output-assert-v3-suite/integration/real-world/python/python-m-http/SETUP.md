# Scenario

**Feature**: python -m http.server

```
# python -m
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"Serving HTTP on :: port 8000 \\(http://\\[::\\]:8000/\\) \\.\\.\\.",
	)
	req.Actual = "Serving HTTP on :: port 8000 (http://[::]:8000/) ..."
	return nil
}
```
