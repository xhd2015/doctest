# Scenario

**Feature**: python -m http.server

```
# python -m
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__LINE__: 'type=string, example=Serving HTTP on :: port 8000 (http://[::]:8000/) ...'\n",
		"__LINE__",
	)
	req.Actual = "Serving HTTP on :: port 8000 (http://[::]:8000/) ..."
	return nil
}
```
