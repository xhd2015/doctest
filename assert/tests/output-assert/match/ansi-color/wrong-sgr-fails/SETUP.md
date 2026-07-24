# Scenario

**Feature**: AC6 — wrong SGR sequence fails

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for AC6 — wrong SGR sequence fails.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "<ansi-color #38;5;208>warn</ansi-color>"
	req.Actual = grayWrap("warn")
	return nil
}
```
