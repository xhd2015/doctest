# Scenario

**Feature**: P16 — block and inline regex

```
# registered tags produce typed AST nodes
Author -> Parser: valid template constructs
Parser -> Pattern with LiteralLine, BlockOptional, Hint, etc.
```

## Steps
1. Set template/actual fields for P16 — block and inline regex.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "<regex>\n^\\.+$\n</regex>\nid=<regex>[0-9]+</regex>"
	return nil
}
```
