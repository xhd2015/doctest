# Scenario

**Bug**: normalization runs only after every successful hook has added its mappings

```
first hook -> shared overlay <- second hook -> active bridge normalizer
```

```go
import (
	"testing"
	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/session"
)
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.PreTest = []core.PreTestHook{
		{Command: []string{"first", "--overlay-file", "$GO_INSTRUMENT_OVERLAY_FILE"}},
		{Command: []string{"second", "--overlay-file", "$GO_INSTRUMENT_OVERLAY_FILE"}},
	}
	req.ActiveBridge, req.CreateActiveBridgeFile = true, true
	req.HookOverlays = [][]OverlayEntry{
		{{Source: "project-source", Replace: "project.go"}},
		{{Source: "active-vendor", Replace: "active.go"}},
	}
	return nil
}
```
