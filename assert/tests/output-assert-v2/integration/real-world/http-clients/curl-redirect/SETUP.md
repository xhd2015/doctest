# Scenario

**Feature**: curl -L

```
# curl redirect
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__LOC__: 'type=string, example=/new'\n",
		"Location: /new",
	)
	req.Actual = "Location: /new"
	return nil
}
```
