# Scenario

**Feature**: pip install

```
# pip install
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__PKG__: 'type=string, example=requests'\n",
		"Successfully installed requests-2\\.31\\.0",
	)
	req.Actual = "Successfully installed requests-2.31.0"
	return nil
}
```
