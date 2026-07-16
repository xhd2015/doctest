# Scenario

**Feature**: gofmt

```
# gofmt -l
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__FILE__: 'type=string, example=main.go'\n",
		"main\\.go",
	)
	req.Actual = "main.go"
	return nil
}
```
