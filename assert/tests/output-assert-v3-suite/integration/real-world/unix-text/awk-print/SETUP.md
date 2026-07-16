# Scenario

**Feature**: awk

```
# awk
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__FIELD__: 'type=string, example=bob'\n",
		"bob",
	)
	req.Actual = "bob"
	return nil
}
```
