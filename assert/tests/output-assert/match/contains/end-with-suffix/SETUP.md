# Scenario

**Feature**: C5 — end-with suffix match

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for C5 — end-with suffix match.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "<contains>\n<end-with>bar</end-with>\n</contains>"
	req.Actual = "line bar"
	return nil
}
```
