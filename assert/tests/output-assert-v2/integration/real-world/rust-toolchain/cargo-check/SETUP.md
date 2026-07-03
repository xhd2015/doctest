# Scenario

**Feature**: cargo check

```
# cargo check
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__LINE__: 'type=string, example=    Finished `dev` profile [unoptimized] target(s) in 0.42s'\n",
		"__LINE__",
	)
	req.Actual = "    Finished `dev` profile [unoptimized] target(s) in 0.42s"
	return nil
}
```
