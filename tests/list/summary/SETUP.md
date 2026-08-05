# Scenario

**Feature**: selection-wide summary after body lines

```
Harness -> list one or more roots
  -> body lines
  -> blank + --- + totals + labels
  -> trailing newline
```

## Preconditions

- Summary always present when ≥1 root listed.
- Totals/labels are sums of per-root body stats.

## Steps

1. Grouping Setup is a no-op.
2. Leaves build multi- or single-root fixtures.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	_ = req
	return nil
}
```
