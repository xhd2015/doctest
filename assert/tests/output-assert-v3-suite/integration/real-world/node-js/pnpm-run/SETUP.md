# Scenario

**Feature**: pnpm run

```
# pnpm run
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__SCRIPT__: 'type=string, example=dev'\n",
		"> dev",
	)
	req.Actual = "> dev"
	return nil
}
```
