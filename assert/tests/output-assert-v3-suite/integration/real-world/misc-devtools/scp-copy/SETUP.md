# Scenario

**Feature**: scp

```
# scp
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__FILE__: 'type=string, example=data.txt'\n",
		"data\\.txt",
	)
	req.Actual = "data.txt"
	return nil
}
```
