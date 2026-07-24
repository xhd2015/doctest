# Scenario

**Feature**: AC4 — raw #90 equals gray

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for AC4 — raw #90 equals gray.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "<ansi-color #90>1 Cached</ansi-color>"
	req.Actual = grayWrap("1 Cached")
	return nil
}
```
