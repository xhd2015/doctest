# Scenario

**Feature**: P14 — raw SGR ansi-color

```
# registered tags produce typed AST nodes
Author -> Parser: valid template constructs
Parser -> Pattern with LiteralLine, BlockOptional, Hint, etc.
```

## Steps
1. Set template/actual fields for P14 — raw SGR ansi-color.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "<ansi-color #38;5;208>x</ansi-color>"
	return nil
}
```
