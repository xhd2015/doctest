# Scenario

**Feature**: A5 — any-of meta lines ignored

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for A5 — any-of meta lines ignored.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "<any-of>\n<expect>\nok\n</expect>\n</any-of>\ndone"
	req.Actual = "ok\ndone"
	return nil
}
```
