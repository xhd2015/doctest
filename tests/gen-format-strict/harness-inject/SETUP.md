# Scenario

**Feature**: harness-injected symbols remain importable without author importing session

```
# minimal author SETUP
import "testing" only
  -> generated leaf still constructs d *session.Doctest
  -> engine may add session / syscall / droot / registry
  -> suite compiles and passes
```

## Preconditions

- User must still import packages **they** name; harness-only injects are allowed.
- FixtureKind `harness-minimal`.

## Steps

1. `Op=run_fixture`, FixtureKind `harness-minimal`.
2. Assert suite success and (optionally) session import present due to harness.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Op = "run_fixture"
	req.FixtureKind = "harness-minimal"
	return nil
}
```
