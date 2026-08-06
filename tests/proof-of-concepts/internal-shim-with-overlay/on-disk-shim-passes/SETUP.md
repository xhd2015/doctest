# Scenario

**Feature**: on-disk bridge under parent of `internal` re-exports the leaf

```
go run ./suite -> imports http/__doctest_internal_shim/leaf -> from-internal-leaf
```

## Steps

1. Set scenario `on-disk-shim`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	_ = d
	req.Scenario = "on-disk-shim"
	return nil
}
```
