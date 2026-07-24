# Scenario

**Feature**: V3-E5 — invalid omit marker syntax

```
# ...abc lines omitted... has non-numeric count
Author -> v3 Parser: bad omit syntax
Parser -> parse error
```

## Steps
1. Set body line `...abc lines omitted...`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template("", "...abc lines omitted...")
	return nil
}
```
