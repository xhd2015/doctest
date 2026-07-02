# Scenario

**Feature**: L4 — trailing blank line in a pure literal template is significant

```
# author writes a pure literal template ending with an empty output line
Author -> Parser: template "foo\n\n"
# strict literal parsing keeps the final blank line
Parser -> Matcher: LiteralLine "foo", LiteralLine "" (trailingNewline = true)
# actual has only the ordinary final newline after foo
Matcher <- actual "foo\n"
Matcher -> fail (trailing blank line missing)
```

## Steps
1. Set template to `"foo\n\n"` (line `foo`, then one trailing blank line).
2. Set actual to `"foo\n"` (no trailing blank line after the `foo` line).

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "foo\n\n"
	req.Actual = "foo\n"
	return nil
}
```
