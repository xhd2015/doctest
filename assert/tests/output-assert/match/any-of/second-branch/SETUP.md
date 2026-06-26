# Scenario

**Feature**: A2 — second any-of branch matches

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for A2 — second any-of branch matches.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "<any-of><expect>a</expect><expect>b</expect></any-of>"
	req.Actual = "b"
	return nil
}
```
