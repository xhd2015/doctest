# Scenario

**Feature**: redis-cli GET

```
# redis-cli GET
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__VAL__: 'type=string, example=hello'\n",
		"hello",
	)
	req.Actual = "hello"
	return nil
}
```
