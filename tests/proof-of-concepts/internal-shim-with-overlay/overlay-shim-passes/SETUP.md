# Scenario

**Feature**: overlay-only virtual bridge under parent of `internal` works

```
go run -overlay=overlay.json ./suite_overlay -> from-internal-leaf
```

## Steps

1. Set scenario `overlay-shim`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	_ = d
	req.Scenario = "overlay-shim"
	return nil
}
```
