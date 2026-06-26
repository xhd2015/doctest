# Scenario

**Feature**: P19 — escaped tag literal in match text

```
# registered tags produce typed AST nodes
Author -> Parser: valid template constructs
Parser -> Pattern with LiteralLine, BlockOptional, Hint, etc.
```

## Steps
1. Set template/actual fields for P19 — escaped tag literal in match text.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "see \\<optional> in docs"
	return nil
}
```
