# Scenario

**Feature**: A1 — missing `import "time"` fails compile; engine does not inject `"time"`

```
fixture leaf SETUP:
  import "testing"
  _ = 20 * time.Second
  -> assemble + WriteFormattedGo
  -> no auto-add of "time"
  -> go test fails (undefined: time)
```

## Preconditions

- Current product still has `stdlibByPkgName` auto-add → expect **RED** until implementer.
- FixtureKind `missing-time`.

## Steps

1. Set `FixtureKind = "missing-time"`.
2. Run generates and attempts suite test.
3. Assert non-empty RunErr and generated sources lack `"time"` import despite `time.Second` usage.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.FixtureKind = "missing-time"
	return nil
}
```
