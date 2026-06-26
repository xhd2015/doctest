# Scenario

**Feature**: Match operation — Pattern vs actual output

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set `req.Operation = "match"`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Operation = "match"
	return nil
}
```
