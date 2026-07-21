# Scenario

**Feature**: user stdlib usage is never auto-injected by the generate write path

```
# A1 missing time
SETUP uses time.Second with import "testing" only
  -> WriteFormattedGo must NOT add "time"
  -> go test fails (undefined: time)

# A2 explicit time
SETUP imports "time" + "testing"
  -> generate + suite go test succeed
```

## Preconditions

- Split factor: presence/absence of explicit `import "time"` when body uses `time.Second`.
- Leaves set `FixtureKind` and share `Op=run_fixture`.

## Steps

1. Set `req.Op = "run_fixture"`.
2. Leaf sets `FixtureKind` to `missing-time` or `explicit-time`.
3. Assert compile fail vs pass and (for A1) no auto-added `"time"` in gen sources.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Op = "run_fixture"
	return nil
}
```
