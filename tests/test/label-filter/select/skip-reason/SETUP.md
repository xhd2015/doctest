# Scenario

**Feature**: label-filter skips carry Reason "label filter"

```
--label manual → every SkippedCase.Reason == "label filter"
```

## Steps

1. Same no-match filter; assert Reason field.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.TreeRoot = writeLabelFilterMod(t)
	req.LabelExprs = []string{"manual"}
	return nil
}
```
