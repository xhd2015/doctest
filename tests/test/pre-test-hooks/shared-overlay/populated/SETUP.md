# Scenario

**Feature**: a populated shared overlay becomes a standard Go argument

```
hook writes overlay -> config driver observes non-empty file -> -overlay=file
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.PreTest = []core.PreTestHook{{Command: []string{"tool", "--overlay-dir", "$GO_INSTRUMENT_OVERLAY_DIR", "--overlay-file", "$GO_INSTRUMENT_OVERLAY_FILE"}}}
	req.WriteOverlay = true
	return nil
}
```
