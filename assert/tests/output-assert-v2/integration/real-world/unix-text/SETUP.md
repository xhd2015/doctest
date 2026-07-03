# Scenario

**Feature**: Unix text and file utilities — v2 CLI templates

```
# unix-text cookbook leaves
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
