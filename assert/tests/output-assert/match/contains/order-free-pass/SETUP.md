# Scenario

**Feature**: C1 — order-free full-line fragments

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for C1 — order-free full-line fragments.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "<contains>\na\nb\n</contains>"
	req.Actual = "z\nb\ny\na"
	return nil
}
```
