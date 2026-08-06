# Scenario

**Feature**: session and cold-cache CLI smoke still exit 0 when env is only via cmd.Env

```
# functional locks
tiny fixture module
  -> doctest test              -> exit 0; nested d.SESSION_ID non-empty
  -> doctest test --cold-cache -> exit 0; cold announce mentions GOCACHE
```

## Preconditions

- Root Setup provides `req.Bin` and `req.ModuleRoot`.
- Leaves create isolated temp modules / cache homes (parallel-safe).
- Full-integration leaves may use **e2e** (nested generate + go test).

## Steps

1. Raise timeout for nested generate + suite go test.
2. Leaves set Args / Env / FixtureDir.

```go
import (
	"testing"
	"time"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Op = "cli"
	req.Timeout = 120 * time.Second
	return nil
}
```
