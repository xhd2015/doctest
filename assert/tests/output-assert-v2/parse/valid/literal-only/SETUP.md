# Scenario

**Feature**: V2-P6 — literal-only body under v2 header

```
# header with only version: 2 + literal lines
Author -> v2 Parser: two literal lines
Parser -> LiteralLine×2
```

## Steps
1. Set body with two literal lines, no placeholders.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template("", "hello\nworld")
	return nil
}
```