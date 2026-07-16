# Scenario

**Feature**: npm init

```
# npm init
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__PATH__: 'type=string, example=/tmp/package.json'\n",
		"Wrote to __PATH__\n...3 lines omitted...\n\\}",
	)
	req.Actual = "Wrote to /tmp/package.json\n  \"name\": \"a\"\n  \"version\": \"1.0.0\"\n  \"main\": \"index.js\"\n}"
	return nil
}
```
