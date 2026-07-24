# Scenario

**Feature**: V3S-E1 — undefined placeholder in body

```
# __MISSING__ used in body but not declared in header
Author -> v3 Parser: undefined placeholder
Parser -> parse error
```

## Steps
1. Set body using `__MISSING__` without header definition.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template("", "Hello __MISSING__")
	return nil
}
```