# Scenario

**Feature**: no user overlay keeps pre_test populated-hook behavior (one `-overlay=`)

```
ApplyPreTestHooksWithUserOverlay(user="", hooks write file)
  -> one -overlay=; Replace from hook only
```

## Preconditions

- `UserOverlayEmpty` true (no user path).
- Hook populates shared overlay file (same intent as
  `tests/test/pre-test-hooks/shared-overlay/populated` — thin regression, not a
  full duplicate suite).

## Steps

1. Empty user overlay.
2. Single hook writes `project-source` → `from-hook`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.UserOverlayEmpty = true
	req.UserReplace = nil
	req.UseMaterializeHelper = false
	req.PreTest = []core.PreTestHook{
		{Command: []string{"tool", "--overlay-file", "$GO_INSTRUMENT_OVERLAY_FILE"}},
	}
	req.HookOverlays = [][]OverlayEntry{
		{{Source: "project-source", Replace: "from-hook"}},
	}
	return nil
}
```
