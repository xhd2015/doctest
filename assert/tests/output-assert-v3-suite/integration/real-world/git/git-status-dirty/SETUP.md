# Scenario

**Feature**: git status dirty

```
# git status
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__FILE__: 'type=string, example=main.go'\n",
		"modified:   main\\.go",
	)
	req.Actual = "modified:   main.go"
	return nil
}
```
