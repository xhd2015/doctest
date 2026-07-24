# Scenario

**Feature**: prettier --check

```
# prettier
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__FILE__: 'type=string, example=src/a.ts'\n",
		"Checking formatting\\.\\.\\.\nsrc/a\\.ts",
	)
	req.Actual = "Checking formatting...\nsrc/a.ts"
	return nil
}
```
