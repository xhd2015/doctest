# Scenario

**Feature**: V3S-P5 — omit marker line

```
# ...3 lines omitted... consumes N lines at match time
Author -> v3 Parser: omit marker template
Parser -> OmitLine{3}
```

## Steps
1. Set body line `...3 lines omitted...`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template("", "...3 lines omitted...")
	return nil
}
```