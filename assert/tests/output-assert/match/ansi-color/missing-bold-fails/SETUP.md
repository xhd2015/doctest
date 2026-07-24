# Scenario

**Feature**: AC10 — missing bold in combo fails

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for AC10 — missing bold in combo fails.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "<ansi-color bold gray>1 Cached</ansi-color>"
	req.Actual = grayWrap("1 Cached")
	return nil
}
```
