# Scenario

**Feature**: P17 — inline any-of pattern line

```
# registered tags produce typed AST nodes
Author -> Parser: valid template constructs
Parser -> Pattern with LiteralLine, BlockOptional, Hint, etc.
```

## Steps
1. Set template/actual fields for P17 — inline any-of pattern line.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "status: <any-of><expect>ok</expect><expect>no</expect></any-of>"
	return nil
}
```
