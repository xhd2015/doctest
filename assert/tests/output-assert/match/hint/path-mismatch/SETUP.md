# Scenario

**Feature**: H4 — hint label in error message

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for H4 — hint label in error message.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "$ cd <hint:path>~/proj</hint:path>"
	req.Actual = "$ cd /tmp/proj"
	return nil
}
```
