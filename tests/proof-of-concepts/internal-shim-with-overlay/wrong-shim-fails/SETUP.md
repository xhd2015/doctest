# Scenario

**Feature**: overlay shim outside parent of `internal` cannot import the leaf

```
go build -overlay=… ./suite_wrong -> use of internal package … not allowed
```

## Steps

1. Set scenario `wrong-shim`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	_ = d
	req.Scenario = "wrong-shim"
	return nil
}
```
