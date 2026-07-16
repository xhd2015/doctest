# Scenario

**Feature**: go vet

```
# go vet
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"main\\.go:10:2: unreachable code",
	)
	req.Actual = "main.go:10:2: unreachable code"
	return nil
}
```
