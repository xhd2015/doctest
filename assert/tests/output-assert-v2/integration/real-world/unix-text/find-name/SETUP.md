# Scenario

**Feature**: find -name

```
# find
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__PATH__: 'type=string, example=./a.go'\n",
		"./a.go",
	)
	req.Actual = "./a.go"
	return nil
}
```
