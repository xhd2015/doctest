# Scenario

**Feature**: P3 — standalone hint segment

```
# registered tags produce typed AST nodes
Author -> Parser: valid template constructs
Parser -> Pattern with LiteralLine, BlockOptional, Hint, etc.
```

## Steps
1. Set template/actual fields for P3 — standalone hint segment.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "<hint:id>abc</hint:id>"
	return nil
}
```
