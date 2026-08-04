# Scenario

**Bug**: an active bridge redirects its matching original vendor source

```
hook -> original vendor key -> active bridge mapping -> bridge source key
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
