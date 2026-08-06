# Scenario

**Feature**: OR expression runs leaves with any matching label

```
--label 'slow || flaky' → run {slow, both, flaky}
```

## Steps

1. Write five-leaf fixture.
2. LabelExprs = ["slow || flaky"].

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.TreeRoot = writeLabelFilterMod(t)
	req.LabelExprs = []string{"slow || flaky"}
	return nil
}
```
