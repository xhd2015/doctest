# Scenario

**Feature**: H2 — hint is not a wildcard

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for H2 — hint is not a wildcard.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "<hint:id>abc</hint:id>"
	req.Actual = "xyz"
	return nil
}
```
