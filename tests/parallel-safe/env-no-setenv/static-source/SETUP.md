# Scenario

**Feature**: product non-test sources must not process-write SESSION_ID or GOCACHE

```
# static scan
moduleRoot/libdoc/**/*.go  (skip *_test.go, testdata/, tests/)
  -> forbid (os|syscall).(Setenv|Unsetenv) of DoctestSessionIDEnv / "DOCTEST_SESSION_ID" / "GOCACHE"
```

## Preconditions

- `req.ModuleRoot` set by root Setup.
- No CLI binary required; `req.Op = static_scan`.

## Steps

1. Mark op as static_scan so Run is a no-op.
2. Assert walks product sources and fails on forbidden writes.

```go
import (
	"testing"
	"time"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Op = "static_scan"
	req.Timeout = 10 * time.Second
	return nil
}
```
