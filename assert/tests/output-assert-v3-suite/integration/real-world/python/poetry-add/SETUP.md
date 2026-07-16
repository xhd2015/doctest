# Scenario

**Feature**: poetry add

```
# poetry add
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__PKG__: 'type=string, example=httpx'\n",
		"Package operations: 1 install, httpx",
	)
	req.Actual = "Package operations: 1 install, httpx"
	return nil
}
```
