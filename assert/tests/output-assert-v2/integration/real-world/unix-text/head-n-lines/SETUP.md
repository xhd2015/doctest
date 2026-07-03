# Scenario

**Feature**: head -n

```
# head file
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"==> f.txt <==\nline one\nline two",
	)
	req.Actual = "==> f.txt <==\nline one\nline two"
	return nil
}
```
