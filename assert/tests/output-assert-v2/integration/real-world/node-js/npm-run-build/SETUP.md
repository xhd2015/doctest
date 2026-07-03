# Scenario

**Feature**: npm run build

```
# npm run
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__BANNER__: 'type=string, example=myapp@1.0.0 build'\n__SEC__: 'type=number, example=1.2'\n",
		"> __BANNER__\n> tsc\n...2 lines omitted...\nDone in __SEC__s.",
	)
	req.Actual = "> myapp@1.0.0 build\n> tsc\nCompiling...\nDone\nDone in 1.2s."
	return nil
}
```
