# Scenario

**Feature**: ripgrep hit

```
# rg PATTERN
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__LINE__: 'type=string, example=src/a.go:1:fmt.Println()'\n",
		"src/a\\.go:1:fmt\\.Println\\(\\)",
	)
	req.Actual = "src/a.go:1:fmt.Println()"
	return nil
}
```
