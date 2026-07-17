# Scenario

**Feature**: filtered suite (`--label expr`) never triggers default-suite warn

```
# LabelExprs non-empty ⇒ not default suite
default_suite=false (label filter), total>0, elapsed>>3m -> false
```

## Preconditions

- Models non-empty LabelExprs as `DefaultSuite=false`.

## Steps

1. Call ShouldWarn with DefaultSuite=false and slow elapsed.
2. Expect false.

## Context

- Same predicate branch as label-all: not default suite.

```go
import (
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/metrics"
)

func Setup(t *testing.T, req *Request) error {
	req.DefaultSuite = false
	req.Total = 2
	req.Elapsed = metrics.DefaultSuiteWarnThreshold + time.Hour
	return nil
}
```
