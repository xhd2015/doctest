# Scenario

**Feature**: docker ps

```
# docker ps
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__ID__: 'type=string, example=abc123'\n",
		"CONTAINER ID   IMAGE\nabc123   nginx",
	)
	req.Actual = "CONTAINER ID   IMAGE\nabc123   nginx"
	return nil
}
```
