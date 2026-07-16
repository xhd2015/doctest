# Scenario

**Feature**: apt install

```
# apt install
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__PKG__: 'type=string, example=curl'\n",
		"Setting up curl \\(8\\.0\\.0\\) \\.\\.\\.",
	)
	req.Actual = "Setting up curl (8.0.0) ..."
	return nil
}
```
