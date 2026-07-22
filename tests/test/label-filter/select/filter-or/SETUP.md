# Scenario

**Feature**: OR expression runs leaves with any matching label

```
--label 'slow || heavy' → run {slow, both, heavy}
```

## Steps

1. Write five-leaf fixture.
2. LabelExprs = ["slow || heavy"].

```go
func Setup(t *testing.T, req *Request) error {
	req.TreeRoot = writeLabelFilterMod(t)
	req.LabelExprs = []string{"slow || heavy"}
	return nil
}
```
