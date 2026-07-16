# Scenario

**Feature**: grep -n

```
# grep -n pattern
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__LINE__: 'type=number, example=3'\n",
		"3:func main\\(\\) \\{",
	)
	req.Actual = "3:func main() {"
	return nil
}
```
