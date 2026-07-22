# Scenario

**Feature**: product CLI help mentions leaf-cache flags (**L3 e2e**, `label: heavy`)

```
doctest test --help
  -> usage includes -a and --no-leaf-cache
```

## Preconditions

- Builds selftest binary (`testbin.Ensure`) — nested product path.
- Leaf labeled `heavy` so default discovery skips it.

## Steps

1. Ensure binary; isolate env not required for help.
2. Child sets Op=`runtime_once` with `test --help`.

```go
import (
	"path/filepath"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/testbin"
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Timeout = 60 * time.Second
	// tests/leaf-cache -> tests -> module root
	modRoot := filepath.Clean(filepath.Join(d.DOCTEST_ROOT, "..", ".."))
	req.Bin = testbin.Ensure(t, modRoot)
	return nil
}
```
