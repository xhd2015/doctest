# Scenario

**Feature**: go doc

```
# go doc fmt
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"package fmt",
	)
	req.Actual = "package fmt"
	return nil
}
```
