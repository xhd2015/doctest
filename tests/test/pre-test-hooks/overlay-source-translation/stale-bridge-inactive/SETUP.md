# Scenario

**Bug**: a stale bridge directory is not current-run bridge metadata

```
hook -> original vendor key; on-disk bridge without mapping -> original key
```

```go
import (
	"testing"
	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/session"
)
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.PreTest = []core.PreTestHook{{Command: []string{"hook", "--overlay-file", "$GO_INSTRUMENT_OVERLAY_FILE"}}}
	req.CreateActiveBridgeFile = true
	req.HookOverlays = [][]OverlayEntry{{{Source: "active-vendor", Replace: "replacement.go"}}}
	return nil
}
```
