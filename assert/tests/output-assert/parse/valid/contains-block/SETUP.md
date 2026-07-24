# Scenario

**Feature**: P9 — contains block with three fragments

```
# registered tags produce typed AST nodes
Author -> Parser: valid template constructs
Parser -> Pattern with LiteralLine, BlockOptional, Hint, etc.
```

## Steps
1. Set template/actual fields for P9 — contains block with three fragments.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "<contains>\nline1\nline2\nline3\n</contains>"
	return nil
}
```
