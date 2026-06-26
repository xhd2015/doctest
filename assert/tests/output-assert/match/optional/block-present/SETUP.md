# Scenario

**Feature**: O2 — block optional present

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for O2 — block optional present.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "head\n<optional>\nnoise\n</optional>\ntail"
	req.Actual = "head\nnoise\ntail"
	return nil
}
```
