# Scenario

**Feature**: Match operation — v2 strict full match

```
# parsed v2 pattern compared to actual output line-by-line
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