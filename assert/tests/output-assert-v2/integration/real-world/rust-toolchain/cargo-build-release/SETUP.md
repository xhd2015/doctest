# Scenario

**Feature**: cargo build --release

```
# cargo build
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__LINE__: 'type=string, example=    Finished `release` profile [optimized] target(s) in 3.21s'\n",
		"__LINE__",
	)
	req.Actual = "    Finished `release` profile [optimized] target(s) in 3.21s"
	return nil
}
```
