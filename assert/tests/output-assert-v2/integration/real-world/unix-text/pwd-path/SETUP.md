# Scenario

**Feature**: pwd

```
# pwd
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__DIR__: 'type=string, example=/tmp/proj'\n",
		"/tmp/proj",
	)
	req.Actual = "/tmp/proj"
	return nil
}
```
