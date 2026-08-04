# Scenario

**Bug**: normalization changes keys only, never project sources or values

```
hook -> project key and vendor key/value -> bridge normalizer -> only vendor key changes
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
	req.HookOverlays = [][]OverlayEntry{{{Source: "project-source", Replace: "replacement.go"}, {Source: "active-vendor", Replace: "replacement.go"}}}
	return nil
}
```
