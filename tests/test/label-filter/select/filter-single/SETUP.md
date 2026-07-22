# Scenario

**Feature**: single-label filter selects all leaves carrying that label

```
--label slow → run {slow, both}; skip {fast, ui, heavy}
```

## Steps

1. Write five-leaf fixture.
2. LabelExprs = ["slow"].

```go
func Setup(t *testing.T, req *Request) error {
	req.TreeRoot = writeLabelFilterMod(t)
	req.LabelExprs = []string{"slow"}
	return nil
}
```
