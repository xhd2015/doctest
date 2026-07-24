# Scenario

**Feature**: repeatable LabelExprs combine with OR (multi-flag)

```
LabelExprs [slow, heavy] ≡ --label slow --label heavy
```

## Steps

1. Write five-leaf fixture.
2. LabelExprs = ["slow", "heavy"].

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.TreeRoot = writeLabelFilterMod(t)
	req.LabelExprs = []string{"slow", "heavy"}
	return nil
}
```
