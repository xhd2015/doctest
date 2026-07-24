# Scenario

**Feature**: O9 — partial inner lines rejected

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for O9 — partial inner lines rejected.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "head\n<optional>\nline1\nline2\n</optional>\ntail"
	req.Actual = "head\nline1\ntail"
	return nil
}
```
