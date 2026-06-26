# Scenario

**Feature**: P8 — block any-of with two branches

```
# registered tags produce typed AST nodes
Author -> Parser: valid template constructs
Parser -> Pattern with LiteralLine, BlockOptional, Hint, etc.
```

## Steps
1. Set template/actual fields for P8 — block any-of with two branches.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "<any-of>\n<expect>\na\n</expect>\n<expect>\nb\n</expect>\n</any-of>"
	return nil
}
```
