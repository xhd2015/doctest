# Scenario

**Feature**: L1 — interior blank line between literals is preserved

```
# author writes literal lines with an interior empty line
Author -> Parser: template "a\n\nb"
# strict parsing preserves the interior empty line
Parser -> Matcher: LiteralLine "a", LiteralLine "", LiteralLine "b"
Matcher <- actual "a\n\nb"
Matcher -> pass
```

## Steps
1. Set template to `"a\n\nb"` (literal `a`, interior empty line, literal `b`).
2. Set actual to `"a\n\nb"` (same bytes, no trailing `\n` on either side).

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "a\n\nb"
	req.Actual = "a\n\nb"
	return nil
}
```
