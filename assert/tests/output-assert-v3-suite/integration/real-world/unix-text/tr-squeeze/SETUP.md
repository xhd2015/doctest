# Scenario

**Feature**: tr -s

```
# tr
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"a b c",
	)
	req.Actual = "a b c"
	return nil
}
```
