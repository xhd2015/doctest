# Scenario

**Feature**: A2 — explicit `import "time"` succeeds

```
fixture leaf SETUP:
  import ("testing"; "time")
  _ = 20 * time.Second
  -> assemble + WriteFormattedGo
  -> go test suite succeeds
```

## Preconditions

- FixtureKind `explicit-time`.
- Same body as A1 but with proper imports.

## Steps

1. Set `FixtureKind = "explicit-time"`.
2. Run generate + suite test.
3. Assert `RunErr` empty.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.FixtureKind = "explicit-time"
	return nil
}
```
