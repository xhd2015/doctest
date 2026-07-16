# Scenario

**Feature**: cargo check

```
# cargo check
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"    Finished `dev` profile \\[unoptimized\\] target\\(s\\) in 0\\.42s",
	)
	req.Actual = "    Finished `dev` profile [unoptimized] target(s) in 0.42s"
	return nil
}
```
