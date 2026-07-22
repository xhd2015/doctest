# Scenario

**Feature**: AND expression requires every listed label on the leaf

```
--label 'slow && ui-automation' → run {both} only
```

## Steps

1. Write five-leaf fixture.
2. LabelExprs = ["slow && ui-automation"].

```go
func Setup(t *testing.T, req *Request) error {
	req.TreeRoot = writeLabelFilterMod(t)
	req.LabelExprs = []string{"slow && ui-automation"}
	return nil
}
```
