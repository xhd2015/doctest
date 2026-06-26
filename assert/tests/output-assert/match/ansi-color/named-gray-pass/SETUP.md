# Scenario

**Feature**: AC1 — named gray wrap passes

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for AC1 — named gray wrap passes.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "<ansi-color gray>1 Cached</ansi-color>"
	req.Actual = grayWrap("1 Cached")
	return nil
}
```
