# Scenario

**Feature**: P11 — named gray ansi-color

```
# registered tags produce typed AST nodes
Author -> Parser: valid template constructs
Parser -> Pattern with LiteralLine, BlockOptional, Hint, etc.
```

## Steps
1. Set template/actual fields for P11 — named gray ansi-color.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "<ansi-color gray>1 Cached</ansi-color>"
	return nil
}
```
