# Scenario

**Feature**: Python ecosystem — v3 CLI templates

```
# python cookbook leaves
Matcher <- simulated tool output
```

## Steps
1. Leaf Setup supplies Template and Actual for one tool transcript.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Operation = "match"
	return nil
}
```
