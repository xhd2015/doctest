# Scenario

**Feature**: warn when default suite is slow and ran tests

```
# default suite, at least one test, elapsed past 3 minutes
default_suite=true, total=1, elapsed=3m+1ns -> ShouldWarn = true
```

## Preconditions

- Threshold is the package default (3 minutes).

## Steps

1. Call ShouldWarn with default suite, total=1, elapsed just over 3m.
2. Expect true.

## Context

- Boundary is strict greater-than: elapsed must exceed threshold.

```go
import (
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/metrics"
)

func Setup(t *testing.T, req *Request) error {
	req.DefaultSuite = true
	req.Total = 1
	req.Elapsed = metrics.DefaultSuiteWarnThreshold + time.Nanosecond
	return nil
}
```
