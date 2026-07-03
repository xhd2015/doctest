# Scenario

**Feature**: docker images

```
# docker images
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__REPO__: 'type=string, example=nginx'\n",
		"nginx              latest",
	)
	req.Actual = "nginx              latest"
	return nil
}
```
