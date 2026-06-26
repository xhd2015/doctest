# Scenario

**Feature**: C3 — default full-line match

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for C3 — default full-line match.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "<contains>\nfoo\n</contains>"
	req.Actual = "x\nxfoo\n"
	return nil
}
```
