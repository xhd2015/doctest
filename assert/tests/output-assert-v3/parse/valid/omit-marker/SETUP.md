# Scenario

**Feature**: V3-P6 — omit marker line

```
# ...3 lines omitted... is special (not content regex)
Author -> v3 Parser: omit marker template
Parser -> OmitLine{3}
```

## Steps
1. Set body line ...3 lines omitted....

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template("", "...3 lines omitted...")
	return nil
}
```
