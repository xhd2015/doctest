# Scenario

**Feature**: single-label filter selects all leaves carrying that label

```
--label slow → run {slow, both}; skip {fast, ui, flaky}
```

## Steps

1. Write five-leaf fixture.
2. LabelExprs = ["slow"].

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.TreeRoot = writeLabelFilterMod(t)
	req.LabelExprs = []string{"slow"}
	return nil
}
```
