# Scenario

**Feature**: V2-M12 — parens without alternation stay pattern line

```
# (1 Cached) has parens but no | alternation — literal
Matcher <- exact (1 Cached)
```

## Steps
1. Set pattern line `(1 Cached)` and identical actual.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template("", "(1 Cached)")
	req.Actual = "(1 Cached)"
	return nil
}
```