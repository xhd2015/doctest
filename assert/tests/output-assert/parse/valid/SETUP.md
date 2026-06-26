# Scenario

**Feature**: Valid parse constructs produce typed AST nodes

```
# registered tags produce typed AST nodes
Author -> Parser: valid template constructs
Parser -> Pattern with LiteralLine, BlockOptional, Hint, etc.
```

## Steps
1. Expect parse success (no `ExpectParseError`).

```go
func Setup(t *testing.T, req *Request) error {
	req.ExpectParseError = false
	return nil
}
```
