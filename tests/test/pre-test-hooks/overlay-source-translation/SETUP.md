# Scenario

**Feature**: post-hook overlay normalization merges only xgo-style phantom go.mod pairs from current-run bridge metadata

```
pre_test hooks -> shared package overlay keys (project vendor)
  + active BridgeRoot placeholders
  -> package keys unchanged; vendor/.../go.mod → placeholder go.mod
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	return nil
}
```
