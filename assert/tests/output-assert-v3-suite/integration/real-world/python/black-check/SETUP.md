# Scenario

**Feature**: black --check

```
# black
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__FILE__: 'type=string, example=main.py'\n",
		"would reformat main\\.py",
	)
	req.Actual = "would reformat main.py"
	return nil
}
```
