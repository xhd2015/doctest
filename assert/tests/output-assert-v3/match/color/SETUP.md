# Scenario

**Feature**: ANSI color spans with QuoteMeta inner text

```
# <ansi-color> structure kept; inner text literal (including dots)
# outside tags remain raw RE; placeholders still expand
Matcher <- ANSI-wrapped actual
```

## Steps
1. Narrow to color match scenarios.

```go
func Setup(t *testing.T, req *Request) error {
	_ = t
	_ = req
	return nil
}
```
