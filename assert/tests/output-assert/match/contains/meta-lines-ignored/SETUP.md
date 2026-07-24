# Scenario

**Feature**: C6 — contains meta lines ignored

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for C6 — contains meta lines ignored.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "<contains>\nfound\n</contains>\nextra"
	req.Actual = "prefix\nfound\nsuffix\nextra"
	return nil
}
```
