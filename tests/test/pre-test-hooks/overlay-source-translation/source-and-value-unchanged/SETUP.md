# Scenario

**Feature**: project and package keys/values stay; only go.mod pair is added

```
hook -> project key + vendor package key (same replacement value)
  + active bridge metadata
  -> keys/values unchanged; go.mod → placeholder merged
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
	req.HookOverlays = [][]OverlayEntry{
		{{Source: "project-source", Replace: "replacement.go"}, {Source: "active-vendor", Replace: "replacement.go"}},
	}
	return nil
}
```
