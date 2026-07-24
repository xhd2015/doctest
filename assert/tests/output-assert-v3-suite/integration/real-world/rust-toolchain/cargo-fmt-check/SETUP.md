# Scenario

**Feature**: cargo fmt --check

```
# cargo fmt
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"Diff in src/lib\\.rs at line 1",
	)
	req.Actual = "Diff in src/lib.rs at line 1"
	return nil
}
```
