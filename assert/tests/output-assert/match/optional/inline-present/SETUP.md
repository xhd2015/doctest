# Scenario

**Feature**: O5 — inline optional present

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for O5 — inline optional present.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "Result: <optional>warn: </optional>OK"
	req.Actual = "Result: warn: OK"
	return nil
}
```
