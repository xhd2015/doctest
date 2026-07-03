# Scenario

**Feature**: gem install

```
# gem install
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__GEM__: 'type=string, example=rails'\n",
		"Successfully installed rails-7.0.0",
	)
	req.Actual = "Successfully installed rails-7.0.0"
	return nil
}
```
