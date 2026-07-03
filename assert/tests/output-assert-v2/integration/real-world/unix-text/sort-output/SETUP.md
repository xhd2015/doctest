# Scenario

**Feature**: sort

```
# sort
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"a\nb\nc",
	)
	req.Actual = "a\nb\nc"
	return nil
}
```
