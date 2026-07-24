# Scenario

**Feature**: M1 — exact single-line match

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for M1 — exact single-line match.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "OK"
	req.Actual = "OK"
	return nil
}
```
