# Scenario

**Feature**: git status dirty

```
# git status
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__FILE__: 'type=string, example=main.go'\n",
		"modified:   main.go",
	)
	req.Actual = "modified:   main.go"
	return nil
}
```
