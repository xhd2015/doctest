# Scenario

**Feature**: A6 — explicit imports compile; do not assert gofmt

```
import ("testing"; "time"); time.Second
  -> suite succeeds
  -> Assert: RunErr empty only (no gofmt check)
```

## Preconditions

- Documents Phase A exit: format.Source not required for a compiling package.

## Steps

1. Inherit FixtureKind `explicit-compile`.
2. Assert compile success without gofmt requirements.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Op = "run_fixture"
	if req.FixtureKind == "" {
		req.FixtureKind = "explicit-compile"
	}
	return nil
}
```
