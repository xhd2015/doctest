# Scenario

**Feature**: X3 — literal prefix required

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for X3 — literal prefix required.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "status: <any-of><expect>ok</expect><expect>err</expect></any-of>"
	req.Actual = "ok"
	return nil
}
```
