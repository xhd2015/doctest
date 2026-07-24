# Scenario

**Feature**: P4 — literal prefix plus hint

```
# registered tags produce typed AST nodes
Author -> Parser: valid template constructs
Parser -> Pattern with LiteralLine, BlockOptional, Hint, etc.
```

## Steps
1. Set template/actual fields for P4 — literal prefix plus hint.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "id=<hint:id>abc</hint:id>"
	return nil
}
```
