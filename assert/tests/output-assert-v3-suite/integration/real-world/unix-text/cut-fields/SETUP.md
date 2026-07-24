# Scenario

**Feature**: cut -f

```
# cut
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__LINE__: 'type=string, example=a\tb'\n",
		"a\tb",
	)
	req.Actual = "a\tb"
	return nil
}
```
