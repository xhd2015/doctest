# Scenario

**Feature**: P1 — two literal lines

```
# registered tags produce typed AST nodes
Author -> Parser: valid template constructs
Parser -> Pattern with LiteralLine, BlockOptional, Hint, etc.
```

## Steps
1. Set template/actual fields for P1 — two literal lines.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "hello\nworld"
	return nil
}
```
