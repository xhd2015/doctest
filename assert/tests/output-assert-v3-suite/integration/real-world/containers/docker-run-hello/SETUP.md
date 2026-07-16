# Scenario

**Feature**: docker run

```
# docker run
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"Hello from Docker!",
	)
	req.Actual = "Hello from Docker!"
	return nil
}
```
