## Preconditions
- Integration harness runs the **doctest binary as a subprocess** (`testbin.Ensure`).
- Stdin is provided via `cmd.Stdin` only. Parent `os.Stdin`/`Stdout`/`Stderr` are never replaced.

## Steps
1. Root setup resolves `req.Bin`.
2. Child SETUP configures Args, StdinSource (`pipe` / `devnull`), and optional Stdin body.
3. Run executes the binary with captured stdout/stderr and cmd.Stdin.

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
