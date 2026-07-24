# Scenario

**Feature**: R1 — block regex matches dot line

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for R1 — block regex matches dot line.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "<regex>\n^\\.+$\n</regex>"
	req.Actual = ".."
	return nil
}
```
