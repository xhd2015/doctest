# Scenario

**Feature**: P2 — block optional between literals

```
# registered tags produce typed AST nodes
Author -> Parser: valid template constructs
Parser -> Pattern with LiteralLine, BlockOptional, Hint, etc.
```

## Steps
1. Set template/actual fields for P2 — block optional between literals.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "before\n<optional>\n</optional>\nafter"
	return nil
}
```
