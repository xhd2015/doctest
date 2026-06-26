# Scenario

**Feature**: O8 — adjacent optionals stay separate

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for O8 — adjacent optionals stay separate.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "before\n<optional>\nline1\n</optional>\n<optional>\nline2\n</optional>\nafter"
	req.Actual = "before\nline2\nafter"
	return nil
}
```
