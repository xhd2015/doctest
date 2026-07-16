# Scenario

**Feature**: Omit marker line consumption

```
# ...N lines omitted... skips N arbitrary actual lines
Matcher consumes fixed N lines between anchors
```

## Steps
1. Templates include omit markers between pattern lines.

```go
func Setup(t *testing.T, req *Request) error {
	_ = t
	_ = req
	return nil
}
```