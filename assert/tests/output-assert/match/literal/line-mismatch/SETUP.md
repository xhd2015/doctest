# Scenario

**Feature**: M3 — line-level mismatch

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for M3 — line-level mismatch.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "a\nb"
	req.Actual = "a\nc"
	return nil
}
```
