# Scenario

**Feature**: A4 — multiline expect branch

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for A4 — multiline expect branch.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "<any-of>\n<expect>\nline1\nline2\n</expect>\n<expect>\nalt\n</expect>\n</any-of>"
	req.Actual = "line1\nline2"
	return nil
}
```
