# Scenario

**Feature**: X1 — inline any-of first branch

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for X1 — inline any-of first branch.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "status: <any-of><expect>ok</expect><expect>err</expect></any-of>"
	req.Actual = "status: ok"
	return nil
}
```
