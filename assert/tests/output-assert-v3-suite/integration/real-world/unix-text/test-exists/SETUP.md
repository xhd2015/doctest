# Scenario

**Feature**: test -f ok

```
# test -f
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"ok",
	)
	req.Actual = "ok"
	return nil
}
```
