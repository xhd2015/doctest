# Scenario

**Feature**: explicit leaf SubDir still honors label filter

```
FilterBySubDir(slow) + LabelExprs → run or skip that leaf only
```

## Steps

1. Write five-leaf fixture; set SubDir to a concrete leaf.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = req
	return nil
}
```
