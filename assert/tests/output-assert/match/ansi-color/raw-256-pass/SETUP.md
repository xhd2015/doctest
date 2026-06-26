# Scenario

**Feature**: AC5 — raw 256-color SGR

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for AC5 — raw 256-color SGR.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "<ansi-color #38;5;208>warn</ansi-color>"
	req.Actual = raw256Wrap("warn")
	return nil
}
```
