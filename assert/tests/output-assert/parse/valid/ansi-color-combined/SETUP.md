# Scenario

**Feature**: P18 — bold gray combined tokens

```
# registered tags produce typed AST nodes
Author -> Parser: valid template constructs
Parser -> Pattern with LiteralLine, BlockOptional, Hint, etc.
```

## Steps
1. Set template/actual fields for P18 — bold gray combined tokens.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "<ansi-color bold gray>x</ansi-color>"
	return nil
}
```
