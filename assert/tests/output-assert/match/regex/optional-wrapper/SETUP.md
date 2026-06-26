# Scenario

**Feature**: R3 — optional regex wrapper absent

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for R3 — optional regex wrapper absent.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "<optional>\n<regex>\n^\\.+$\n</regex>\n</optional>\n  (2 Run, 2 Pass, 0 Cached, 0 Fail)"
	req.Actual = "  (2 Run, 2 Pass, 0 Cached, 0 Fail)"
	return nil
}
```
