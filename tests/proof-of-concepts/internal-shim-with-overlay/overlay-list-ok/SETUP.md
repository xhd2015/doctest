# Scenario

**Feature**: `go list -overlay` sees the virtual shim package

```
go list -overlay=overlay.json example.com/realmod/http/__doctest_internal_shim_overlay/leaf
```

## Steps

1. Set scenario `overlay-list`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	_ = d
	req.Scenario = "overlay-list"
	return nil
}
```
