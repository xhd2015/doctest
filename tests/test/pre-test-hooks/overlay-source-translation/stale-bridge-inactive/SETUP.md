# Scenario

**Feature**: on-disk placeholder without current-run metadata does not merge

```
hook -> package key; CreateActiveBridgeFile without ActiveBridge -> no go.mod merge
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.PreTest = []core.PreTestHook{{Command: []string{"hook", "--overlay-file", "$GO_INSTRUMENT_OVERLAY_FILE"}}}
	req.CreateActiveBridgeFile = true // placeholder exists on disk but bridges slice empty
	req.HookOverlays = [][]OverlayEntry{{{Source: "active-vendor", Replace: "replacement.go"}}}
	return nil
}
```
