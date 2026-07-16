# Scenario

**Feature**: cargo init

```
# cargo init
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__NAME__: 'type=string, example=myapp'\n",
		"     Created binary \\(application\\) `myapp` package",
	)
	req.Actual = "     Created binary (application) `myapp` package"
	return nil
}
```
