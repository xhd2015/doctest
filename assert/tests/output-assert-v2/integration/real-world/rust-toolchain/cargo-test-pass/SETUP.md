# Scenario

**Feature**: cargo test

```
# cargo test
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__N__: 'type=number, example=3'\n",
		"test result: ok. 3 passed; 0 failed",
	)
	req.Actual = "test result: ok. 3 passed; 0 failed"
	return nil
}
```
