# Scenario

**Feature**: cat file dump

```
# cat README.md
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"# My Project\nVersion 1\\.0",
	)
	req.Actual = "# My Project\nVersion 1.0"
	return nil
}
```
