# Scenario

**Feature**: Node.js / npm ecosystem — v3 CLI templates

```
# node-js cookbook leaves
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
