# Scenario

**Feature**: uniq -c

```
# uniq -c
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__LINE__: 'type=string, example=   3 foo'\n",
		"   3 foo",
	)
	req.Actual = "   3 foo"
	return nil
}
```
