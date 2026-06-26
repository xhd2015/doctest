# Scenario

**Feature**: M2 — strict trailing newline mismatch

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for M2 — strict trailing newline mismatch.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "OK"
	req.Actual = "OK\n"
	return nil
}
```
