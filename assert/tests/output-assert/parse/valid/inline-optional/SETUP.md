# Scenario

**Feature**: P7 — inline optional segment

```
# registered tags produce typed AST nodes
Author -> Parser: valid template constructs
Parser -> Pattern with LiteralLine, BlockOptional, Hint, etc.
```

## Steps
1. Set template/actual fields for P7 — inline optional segment.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "Result: <optional>warn: </optional>OK"
	return nil
}
```
