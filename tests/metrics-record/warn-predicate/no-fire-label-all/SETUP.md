# Scenario

**Feature**: labeled full suite (`--label-all`) never triggers default-suite warn

```
# not default suite even when very slow with tests
default_suite=false (label-all), total>0, elapsed>>3m -> false
```

## Preconditions

- Models `LabelAll=true` as `DefaultSuite=false`.

## Steps

1. Call ShouldWarn with DefaultSuite=false, total=5, elapsed well over threshold.
2. Expect false.

## Context

- Default suite is only when `!LabelAll && len(LabelExprs)==0`.

```go
import (
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/metrics"
)

func Setup(t *testing.T, req *Request) error {
	req.DefaultSuite = false
	req.Total = 5
	req.Elapsed = metrics.DefaultSuiteWarnThreshold + 10*time.Minute
	return nil
}
```
