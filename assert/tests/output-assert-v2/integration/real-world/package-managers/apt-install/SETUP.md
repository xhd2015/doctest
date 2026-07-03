# Scenario

**Feature**: apt install

```
# apt install
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__LINE__: 'type=string, example=Setting up curl (8.0.0) ...'\n",
		"__LINE__",
	)
	req.Actual = "Setting up curl (8.0.0) ..."
	return nil
}
```
