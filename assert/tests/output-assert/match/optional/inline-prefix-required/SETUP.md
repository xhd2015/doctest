# Scenario

**Feature**: O6 — literal prefix required

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for O6 — literal prefix required.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "Result: <optional>warn: </optional>OK"
	req.Actual = "OK"
	return nil
}
```
