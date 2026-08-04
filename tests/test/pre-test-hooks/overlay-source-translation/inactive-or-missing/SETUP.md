# Scenario

**Bug**: a mapping must match the module and an existing bridged source

```
hook -> inactive and missing-active vendor keys -> normalizer -> original keys
```

```go
import (
	"testing"
	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/session"
)
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.PreTest = []core.PreTestHook{{Command: []string{"hook", "--overlay-file", "$GO_INSTRUMENT_OVERLAY_FILE"}}}
	req.ActiveBridge = true
	req.HookOverlays = [][]OverlayEntry{{{Source: "active-vendor", Replace: "missing.go"}, {Source: "inactive-vendor", Replace: "inactive.go"}}}
	return nil
}
```
