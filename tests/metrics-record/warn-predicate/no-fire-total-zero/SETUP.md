# Scenario

**Feature**: empty suite (total=0) never warns even if wall clock is large

```
# default suite but no tests counted
default_suite=true, total=0, elapsed>>3m -> false
```

## Preconditions

- Covers discovery that finds nothing runnable / all skipped without total>0.

## Steps

1. Call ShouldWarn with total=0 and large elapsed.
2. Expect false.

## Context

- Requirement: `total > 0` is required for WARNING.

```go
import (
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/metrics"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DefaultSuite = true
	req.Total = 0
	req.Elapsed = metrics.DefaultSuiteWarnThreshold + time.Hour
	return nil
}
```
