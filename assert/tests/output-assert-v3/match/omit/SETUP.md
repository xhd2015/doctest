# Scenario

**Feature**: Omit marker consumes exactly N actual lines

```
# ...N lines omitted... is special (not content regex)
Matcher skips N middle lines between anchors
```

## Steps
1. Narrow to omit match scenarios.

```go
func Setup(t *testing.T, req *Request) error {
	_ = t
	_ = req
	return nil
}
```
