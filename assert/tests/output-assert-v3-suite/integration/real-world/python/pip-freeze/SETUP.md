# Scenario

**Feature**: pip freeze

```
# pip freeze
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"requests==2\\.31\\.0",
	)
	req.Actual = "requests==2.31.0"
	return nil
}
```
