# Scenario

**Feature**: Parse operation — template to AST

```
# template text becomes AST
Author -> Parser: template
Parser -> Pattern summary (parse-only leaves)
```

## Steps
1. Set `req.Operation = "parse"`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Operation = "parse"
	return nil
}
```
