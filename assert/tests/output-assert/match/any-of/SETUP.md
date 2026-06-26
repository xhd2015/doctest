# Scenario

**Feature**: Any-of branch selection

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Templates use `<any-of>`.

```go
func Setup(t *testing.T, req *Request) error {
	_ = t
	return nil
}
```
