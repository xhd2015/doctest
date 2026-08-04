# Scenario

**Feature**: mid-string overlay file placeholder allocates empty JSON and substitutes

```
arg "--overlay=$GO_INSTRUMENT_OVERLAY_FILE" -> driver allocates file -> executor receives "--overlay=<abs path>"
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.PreTest = []core.PreTestHook{{Command: []string{"tool", "--overlay=$GO_INSTRUMENT_OVERLAY_FILE"}}}
	return nil
}
```
