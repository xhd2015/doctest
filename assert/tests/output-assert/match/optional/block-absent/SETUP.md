# Scenario

**Feature**: O1 — block optional absent

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for O1 — block optional absent.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "head\n<optional>\n</optional>\ntail"
	req.Actual = "head\ntail"
	return nil
}
```
