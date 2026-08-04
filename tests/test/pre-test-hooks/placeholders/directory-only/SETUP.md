# Scenario

**Feature**: directory placeholder allocates the shared generated directory only

```
directory placeholder -> config driver -> hook receives generated directory
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.PreTest = []core.PreTestHook{{Command: []string{"tool", "--overlay-dir", "$GO_INSTRUMENT_OVERLAY_DIR"}}}
	return nil
}
```
