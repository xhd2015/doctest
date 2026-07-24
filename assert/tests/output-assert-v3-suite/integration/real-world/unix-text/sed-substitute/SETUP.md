# Scenario

**Feature**: sed

```
# sed s///
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"hello earth",
	)
	req.Actual = "hello earth"
	return nil
}
```
