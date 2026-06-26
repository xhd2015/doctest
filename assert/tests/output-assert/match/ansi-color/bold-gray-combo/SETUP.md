# Scenario

**Feature**: AC9 — bold gray strict order

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for AC9 — bold gray strict order.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "<ansi-color bold gray>1 Cached</ansi-color>"
	req.Actual = boldGrayWrap("1 Cached")
	return nil
}
```
