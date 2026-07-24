# Scenario

**Feature**: expression matching no leaves yields empty run set

```
--label manual → run {}; skip all five
```

## Steps

1. Write five-leaf fixture.
2. LabelExprs = ["manual"].

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.TreeRoot = writeLabelFilterMod(t)
	req.LabelExprs = []string{"manual"}
	return nil
}
```
