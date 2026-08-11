# Scenario

**Feature**: multi-leaf parent-internal supports single-package `-coverprofile`

```
RunTest(multi-leaf parent internal, Cover, CoverProfile=…/cover.out)
  -> exit 0 (single suite package)
  -> cover.out exists non-empty
```

## Preconditions

- Fixture from parent multi-leaf Setup.
- Cover profile path is absolute under `t.TempDir()` (parallel-safe).

## Steps

1. Set `WithCover` and absolute `CoverPath`.
2. Run subject tree; Assert profile file + success.

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	_ = d
	req.WithCover = true
	req.CoverPath = filepath.Join(t.TempDir(), "cover.out")
	_ = os.Remove(req.CoverPath)
	return nil
}
```
