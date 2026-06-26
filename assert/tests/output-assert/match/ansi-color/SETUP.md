# Scenario

**Feature**: ANSI color span assertions

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Templates use `<ansi-color>`.

```go
func Setup(t *testing.T, req *Request) error {
	_ = t
	return nil
}
```
