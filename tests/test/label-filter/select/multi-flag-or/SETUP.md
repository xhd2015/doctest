# Scenario

**Feature**: repeatable LabelExprs combine with OR (multi-flag)

```
LabelExprs [slow, flaky] ≡ --label slow --label flaky
```

## Steps

1. Write five-leaf fixture.
2. LabelExprs = ["slow", "flaky"].

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.TreeRoot = writeLabelFilterMod(t)
	req.LabelExprs = []string{"slow", "flaky"}
	return nil
}
```
