# Scenario

**Feature**: AC2 — plain text without color fails

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for AC2 — plain text without color fails.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "<ansi-color gray>1 Cached</ansi-color>"
	req.Actual = "1 Cached"
	return nil
}
```
