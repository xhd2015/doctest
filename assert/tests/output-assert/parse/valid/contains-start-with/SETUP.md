# Scenario

**Feature**: P12 — start-with fragment inside contains

```
# registered tags produce typed AST nodes
Author -> Parser: valid template constructs
Parser -> Pattern with LiteralLine, BlockOptional, Hint, etc.
```

## Steps
1. Set template/actual fields for P12 — start-with fragment inside contains.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "<contains>\n<start-with>\n  agent\n</start-with>\n</contains>"
	return nil
}
```
