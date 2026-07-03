# Scenario

**Feature**: cut -f

```
# cut
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__LINE__: 'type=string, example=a\tb'\n",
		"a\tb",
	)
	req.Actual = "a\tb"
	return nil
}
```
