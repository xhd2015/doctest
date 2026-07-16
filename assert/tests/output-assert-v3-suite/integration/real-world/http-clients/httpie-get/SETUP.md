# Scenario

**Feature**: http GET

```
# http
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__CODE__: 'type=number, example=200'\n",
		"HTTP/1\\.1 200 OK",
	)
	req.Actual = "HTTP/1.1 200 OK"
	return nil
}
```
