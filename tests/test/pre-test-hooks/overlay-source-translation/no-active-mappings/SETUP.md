# Scenario

**Bug**: a no-vendor run has no active mappings to translate

```
hook -> vendor keys -> empty current-run bridge map -> original keys
```

```go
import (
	"testing"
	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/session"
)
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.PreTest = []core.PreTestHook{{Command: []string{"hook", "--overlay-file", "$GO_INSTRUMENT_OVERLAY_FILE"}}}
	req.HookOverlays = [][]OverlayEntry{{{Source: "active-vendor", Replace: "active.go"}, {Source: "inactive-vendor", Replace: "inactive.go"}}}
	return nil
}
```
