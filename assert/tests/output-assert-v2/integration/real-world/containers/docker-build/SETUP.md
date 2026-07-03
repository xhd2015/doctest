# Scenario

**Feature**: docker build

```
# docker build
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__TAG__: 'type=string, example=myimg:latest'\n",
		"Successfully tagged myimg:latest",
	)
	req.Actual = "Successfully tagged myimg:latest"
	return nil
}
```
