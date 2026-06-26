# Scenario

**Feature**: AC3 — named green passes

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for AC3 — named green passes.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "<ansi-color green>2 Pass</ansi-color>"
	req.Actual = greenWrap("2 Pass")
	return nil
}
```
