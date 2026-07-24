# Scenario

**Feature**: Parse operation — v3 template to AST

```
# template text becomes AST via Facade.Parse
Author -> Facade: template
Facade -> v3 Parser or legacy_v1
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