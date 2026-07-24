# Scenario

**Feature**: Literal line exact matching

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Templates are plain literal lines.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	return nil
}
```
