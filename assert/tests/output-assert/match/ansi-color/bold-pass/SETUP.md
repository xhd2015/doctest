# Scenario

**Feature**: AC8 — bold token passes

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for AC8 — bold token passes.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "<ansi-color bold>Hi</ansi-color>"
	req.Actual = boldWrap("Hi")
	return nil
}
```
