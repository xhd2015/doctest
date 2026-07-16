# Scenario

**Feature**: pwd

```
# pwd
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__DIR__: 'type=string, example=/tmp/proj'\n",
		"/tmp/proj",
	)
	req.Actual = "/tmp/proj"
	return nil
}
```
