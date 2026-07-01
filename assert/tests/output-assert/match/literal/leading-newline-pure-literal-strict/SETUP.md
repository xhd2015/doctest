# Scenario

**Feature**: L2 — pure-literal with leading \n still enforces strict trailing policy

```
# author writes a single literal with a leading \n (no trailing \n)
Author -> Parser: template "\nfoo"
# P4 strips the leading empty line, leaving a pure single-literal "foo"
Parser -> Matcher: LiteralLine "foo" (trailingNewline = false)
# strict trailing policy still applies: template has no \n, actual has \n
Matcher <- actual "foo\n"
Matcher -> fail (trailing newline policy violated)
```

## Steps
1. Set template to `"\nfoo"` (leading `\n` then single literal `foo`, no trailing `\n`).
2. Set actual to `"foo\n"` (ends with `\n`).

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "\nfoo"
	req.Actual = "foo\n"
	return nil
}
```
