# Scenario

**Feature**: wget

```
# wget
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__FILE__: 'type=string, example=file.zip'\n",
		"Saving to: '__FILE__",
	)
	req.Actual = "Saving to: 'file.zip'"
	return nil
}
```
