# Scenario

**Feature**: format write path does not auto-inject session for bare user `session.` usage

```
# synthetic package (no assemble harness)
source references session.Once without import
  -> core.WriteFormattedGo
  -> written file still lacks session import
  -> go build fails (undefined: session)
```

## Preconditions

- This isolates **format reconcile** from assemble harness force-add of session.
- Optional-if-hard for full assemble path; library WriteFormattedGo is the strict signal.
- Op `write_format_build`.

## Steps

1. Leaf supplies Source with `session.` selector and no session import.
2. Run WriteFormattedGo + go test.
3. Assert no session import in written source and build fails.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Op = "write_format_build"
	req.WantBuild = true
	return nil
}
```
