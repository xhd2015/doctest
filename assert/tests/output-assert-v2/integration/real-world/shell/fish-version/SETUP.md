# Scenario

**Feature**: fish --version

```
# fish
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"fish, version 3.6.0",
	)
	req.Actual = "fish, version 3.6.0"
	return nil
}
```
