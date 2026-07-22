## Preconditions
- Integration harness runs the **doctest binary as a subprocess** (`testbin.Ensure`).
- Scenario env is applied only via `cmd.Env` (child process). Parent process env is never mutated.

## Steps
1. Root setup resolves `req.Bin` via `testbin.Ensure`.
2. Child SETUP.md files configure `Args` and optional `Env` (KEY=VAL for the child only).
3. Run executes `req.Bin` with `req.Args`, `cmd.Env = merge(parentEnviron, req.Env)`, captures stdout/stderr/exit.

```go
import (
	"path/filepath"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/testbin"
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Timeout == 0 {
		req.Timeout = 60 * time.Second
	}
	if req.Bin == "" {
		mod := filepath.Join(d.DOCTEST_ROOT, "..", "..", "..", "..")
		req.Bin = testbin.Ensure(t, mod)
	}
	return nil
}
```
