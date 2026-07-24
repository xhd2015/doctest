# Scenario

**Feature**: cargo build --release

```
# cargo build
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"    Finished `release` profile \\[optimized\\] target\\(s\\) in 3\\.21s",
	)
	req.Actual = "    Finished `release` profile [optimized] target(s) in 3.21s"
	return nil
}
```
