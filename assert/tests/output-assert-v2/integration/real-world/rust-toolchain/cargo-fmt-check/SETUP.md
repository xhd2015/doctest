# Scenario

**Feature**: cargo fmt --check

```
# cargo fmt
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"Diff in src/lib.rs at line 1",
	)
	req.Actual = "Diff in src/lib.rs at line 1"
	return nil
}
```
