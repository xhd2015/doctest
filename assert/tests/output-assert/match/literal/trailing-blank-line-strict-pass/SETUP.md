# Scenario

**Feature**: L5 — actual output with the same trailing blank line matches

```
# author writes a pure literal template ending with an empty output line
Author -> Parser: template "foo\n\n"
# strict literal parsing keeps the final blank line
Parser -> Matcher: LiteralLine "foo", LiteralLine "" (trailingNewline = true)
# actual includes the same trailing blank line
Matcher <- actual "foo\n\n"
Matcher -> pass
```

## Steps
1. Set template to `"foo\n\n"`.
2. Set actual to `"foo\n\n"`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "foo\n\n"
	req.Actual = "foo\n\n"
	return nil
}
```
