# Scenario

**Feature**: sort

```
# sort
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"a\nb\nc",
	)
	req.Actual = "a\nb\nc"
	return nil
}
```
