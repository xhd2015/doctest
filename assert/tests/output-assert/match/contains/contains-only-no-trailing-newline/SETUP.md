# Scenario

**Feature**: C8 — contains-only pattern ignores trailing newline policy

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for C8 — contains-only pattern ignores trailing newline policy.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "<contains>\nfoo\n</contains>"
	req.Actual = "bar\nfoo\n"
	return nil
}
```
