# Scenario

**Feature**: L2 — leading newline in a pure literal template is significant

```
# author writes a pure literal template with a leading \n
Author -> Parser: template "\nfoo"
# strict literal parsing keeps the leading empty line
Parser -> Matcher: LiteralLine "", LiteralLine "foo" (trailingNewline = false)
# actual omits the leading blank line
Matcher <- actual "foo"
Matcher -> fail (leading blank line missing)
```

## Steps
1. Set template to `"\nfoo"` (leading `\n` then single literal `foo`, no trailing `\n`).
2. Set actual to `"foo"` (same text without the leading blank line).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "\nfoo"
	req.Actual = "foo"
	return nil
}
```
