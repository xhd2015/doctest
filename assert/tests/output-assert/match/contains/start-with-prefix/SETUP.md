# Scenario

**Feature**: C4 — start-with prefix match

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for C4 — start-with prefix match.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "<contains>\n<start-with>foo</start-with>\n</contains>"
	req.Actual = "x\nxfoo extra"
	return nil
}
```
