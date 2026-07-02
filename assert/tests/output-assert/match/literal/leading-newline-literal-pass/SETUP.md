# Scenario

**Feature**: L3 — actual output with the same leading newline matches

```
# author writes a pure literal template with a leading \n
Author -> Parser: template "\nfoo"
# strict literal parsing keeps the leading empty line
Parser -> Matcher: LiteralLine "", LiteralLine "foo" (trailingNewline = false)
# actual includes the same leading blank line
Matcher <- actual "\nfoo"
Matcher -> pass
```

## Steps
1. Set template to `"\nfoo"`.
2. Set actual to `"\nfoo"`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "\nfoo"
	req.Actual = "\nfoo"
	return nil
}
```
