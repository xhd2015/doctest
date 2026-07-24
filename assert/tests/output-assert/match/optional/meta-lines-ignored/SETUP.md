# Scenario

**Feature**: O7 — block meta lines never compared

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for O7 — block meta lines never compared.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "before\n<optional>\ninner\n</optional>\nafter"
	req.Actual = "before\ninner\nafter"
	return nil
}
```
