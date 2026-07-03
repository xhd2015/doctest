# Scenario

**Feature**: scp

```
# scp
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__FILE__: 'type=string, example=data.txt'\n",
		"data.txt",
	)
	req.Actual = "data.txt"
	return nil
}
```
