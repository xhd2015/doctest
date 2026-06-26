# Scenario

**Feature**: A3 — no branch matches, all reported

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for A3 — no branch matches, all reported.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "<any-of><expect>a</expect><expect>b</expect></any-of>"
	req.Actual = "c"
	return nil
}
```
