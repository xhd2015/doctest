# Scenario

**Feature**: git tag

```
# git tag
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__TAG__: 'type=string, example=v1.0.0'\n",
		"v1\\.0\\.0",
	)
	req.Actual = "v1.0.0"
	return nil
}
```
