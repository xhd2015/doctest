# Scenario

**Feature**: comm -13

```
# comm
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"only-b",
	)
	req.Actual = "only-b"
	return nil
}
```
