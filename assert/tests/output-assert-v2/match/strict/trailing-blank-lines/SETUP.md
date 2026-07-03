# Scenario

**Feature**: V2-M17 — trailing blank body lines are not trimmed

```
# body ends with blank line before optional trailing newline
Matcher preserves interior/trailing blank lines in template body
```

## Steps
1. Template body ends with a blank line; actual matches including that blank line.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = `---
version: 2
---
first

`
	req.Actual = "first\n\n"
	return nil
}
```