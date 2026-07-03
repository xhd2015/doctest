# Scenario

**Feature**: go vet

```
# go vet
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"main.go:10:2: unreachable code",
	)
	req.Actual = "main.go:10:2: unreachable code"
	return nil
}
```
