# Scenario

**Feature**: V2-E1 — undefined placeholder in body

```
# __MISSING__ used in body but not declared in header
Author -> v2 Parser: undefined placeholder
Parser -> parse error
```

## Steps
1. Set body using `__MISSING__` without header definition.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template("", "Hello __MISSING__")
	return nil
}
```