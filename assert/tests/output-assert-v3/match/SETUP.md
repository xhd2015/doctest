# Scenario

**Feature**: Match operation — v3 strict full match

```
# parsed v3 pattern compared to actual output line-by-line
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered / binding error
```

## Steps
1. Set `req.Operation = "match"`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Operation = "match"
	return nil
}
```
