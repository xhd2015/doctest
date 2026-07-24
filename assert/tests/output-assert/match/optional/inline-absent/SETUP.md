# Scenario

**Feature**: O4 — inline optional absent

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for O4 — inline optional absent.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "Result: <optional>warn: </optional>OK"
	req.Actual = "Result: OK"
	return nil
}
```
