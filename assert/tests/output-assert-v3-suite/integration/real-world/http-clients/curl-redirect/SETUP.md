# Scenario

**Feature**: curl -L

```
# curl redirect
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__LOC__: 'type=string, example=/new'\n",
		"Location: /new",
	)
	req.Actual = "Location: /new"
	return nil
}
```
