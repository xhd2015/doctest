# Scenario

**Feature**: C7 — contains fragment with inline any-of

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for C7 — contains fragment with inline any-of.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "<contains>\ndefault\n<any-of><expect>not configured</expect><expect>profiles</expect><expect>no profiles</expect></any-of>\n</contains>"
	req.Actual = "default\nnot configured\nno profiles"
	return nil
}
```
