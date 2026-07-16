# Scenario

**Feature**: go build fail

```
# go build
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__PKG__: 'type=string, example=example.com/x'\n",
		"# __PKG__\n\\./main\\.go:2: undefined: X\n...1 lines omitted...\nFAIL",
	)
	req.Actual = "# example.com/x\n./main.go:2: undefined: X\nnote: see docs\nFAIL"
	return nil
}
```
