# Scenario

**Feature**: Go toolchain — v3 CLI templates

```
# go-toolchain cookbook leaves
Matcher <- simulated tool output
```

## Steps
1. Leaf Setup supplies Template and Actual for one tool transcript.

```go
func Setup(t *testing.T, req *Request) error {
	req.Operation = "match"
	return nil
}
```
