# Scenario

**Feature**: active xgo-style bridge merges phantom go.mod; package key stays on project vendor

```
hook -> package overlay on vendor/.../active.go
  + active BridgeRoot = vendor-gomod-overlay/.../go.mod
  -> package key unchanged; go.mod → placeholder merged
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.PreTest = []core.PreTestHook{{Command: []string{"hook", "--overlay-file", "$GO_INSTRUMENT_OVERLAY_FILE"}}}
	req.ActiveBridge, req.CreateActiveBridgeFile = true, true
	req.HookOverlays = [][]OverlayEntry{{{Source: "active-vendor", Replace: "replacement.go"}}}
	return nil
}
```
