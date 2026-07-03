# Scenario

**Feature**: gofmt

```
# gofmt -l
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__FILE__: 'type=string, example=main.go'\n",
		"main.go",
	)
	req.Actual = "main.go"
	return nil
}
```
