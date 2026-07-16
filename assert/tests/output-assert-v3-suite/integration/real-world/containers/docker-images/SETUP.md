# Scenario

**Feature**: docker images

```
# docker images
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__REPO__: 'type=string, example=nginx'\n",
		"nginx              latest",
	)
	req.Actual = "nginx              latest"
	return nil
}
```
