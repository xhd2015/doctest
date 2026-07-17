# Scenario

**Feature**: no warn when default suite finishes within the 3-minute budget

```
# at threshold and under threshold both false
elapsed == 3m -> false
elapsed < 3m  -> false
```

## Preconditions

- Default suite with total > 0.

## Steps

1. Evaluate two cases via `req.WarnCases`: elapsed == threshold and elapsed < threshold.
2. Both must be false.

## Context

- Strict `>` comparison; equality does not warn.

```go
import (
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/metrics"
)

func Setup(t *testing.T, req *Request) error {
	th := metrics.DefaultSuiteWarnThreshold
	req.WarnCases = []WarnCase{
		{
			Name:         "at-threshold",
			DefaultSuite: true,
			Total:        3,
			Elapsed:      th,
			Want:         false,
		},
		{
			Name:         "under-threshold",
			DefaultSuite: true,
			Total:        3,
			Elapsed:      th - time.Second,
			Want:         false,
		},
	}
	return nil
}
```
