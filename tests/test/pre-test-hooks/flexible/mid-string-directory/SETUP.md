# Scenario

**Feature**: mid-string overlay directory placeholder allocates and keeps path suffix

```
arg "--dir=$GO_INSTRUMENT_OVERLAY_DIR/extra" -> driver allocates dir -> executor receives "--dir=<abs>/extra"
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.PreTest = []core.PreTestHook{{Command: []string{"tool", "--dir=$GO_INSTRUMENT_OVERLAY_DIR/extra"}}}
	return nil
}
```
