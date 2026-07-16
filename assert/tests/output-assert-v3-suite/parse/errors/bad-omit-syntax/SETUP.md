# Scenario

**Feature**: V3S-E3 — invalid omit marker syntax

```
# ...abc lines omitted... has non-numeric count
Author -> v3 Parser: bad omit syntax
Parser -> parse error
```

## Steps
1. Set body line `...abc lines omitted...` (omit-intent shape, non-numeric count — do not escape the dots; omit is special, not content RE).

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template("", "...abc lines omitted...")
	return nil
}
```
