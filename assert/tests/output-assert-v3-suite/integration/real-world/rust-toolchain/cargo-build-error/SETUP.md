# Scenario

**Feature**: cargo build error

```
# cargo build
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"error\\[E0425\\]: cannot find value `X` in this scope",
	)
	req.Actual = "error[E0425]: cannot find value `X` in this scope"
	return nil
}
```
