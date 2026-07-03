# Scenario

**Feature**: prettier --check

```
# prettier
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__FILE__: 'type=string, example=src/a.ts'\n",
		"Checking formatting...\nsrc/a.ts",
	)
	req.Actual = "Checking formatting...\nsrc/a.ts"
	return nil
}
```
