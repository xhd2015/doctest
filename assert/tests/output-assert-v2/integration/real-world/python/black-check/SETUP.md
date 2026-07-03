# Scenario

**Feature**: black --check

```
# black
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__FILE__: 'type=string, example=main.py'\n",
		"would reformat main.py",
	)
	req.Actual = "would reformat main.py"
	return nil
}
```
