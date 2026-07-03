# Scenario

**Feature**: V2-E3 — invalid omit marker syntax

```
# ...abc lines omitted... has non-numeric count
Author -> v2 Parser: bad omit syntax
Parser -> parse error
```

## Steps
1. Set body line `...abc lines omitted...`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template("", "...abc lines omitted...")
	return nil
}
```