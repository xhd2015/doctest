# Scenario

**Feature**: M2b — trailing newline agreement

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for M2b — trailing newline agreement.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "OK\n"
	req.Actual = "OK\n"
	return nil
}
```
