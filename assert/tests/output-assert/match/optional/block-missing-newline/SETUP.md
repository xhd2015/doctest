# Scenario

**Feature**: O3 — missing newline between anchors

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for O3 — missing newline between anchors.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "head\n<optional>\n</optional>\ntail"
	req.Actual = "headtail"
	return nil
}
```
