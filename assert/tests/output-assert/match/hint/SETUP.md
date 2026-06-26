# Scenario

**Feature**: Hint literal matching

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Templates use `<hint:label>`.

```go
func Setup(t *testing.T, req *Request) error {
	_ = t
	return nil
}
```
