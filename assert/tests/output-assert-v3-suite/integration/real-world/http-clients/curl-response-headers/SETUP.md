# Scenario

**Feature**: curl -i

```
# curl -i
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__CODE__: 'type=number, example=200'\n",
		"HTTP/1\\.1 __CODE__ OK\nContent-Type: application/json\n...2 lines omitted...\n\\{\"ok\":true\\}",
	)
	req.Actual = "HTTP/1.1 200 OK\nContent-Type: application/json\nContent-Length: 11\n\n{\"ok\":true}"
	return nil
}
```
