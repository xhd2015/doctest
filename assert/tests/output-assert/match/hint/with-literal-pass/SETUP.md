# Scenario

**Feature**: H3 — literal prefix with hint

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for H3 — literal prefix with hint.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "id=<hint:id>abc</hint:id>"
	req.Actual = "id=abc"
	return nil
}
```
