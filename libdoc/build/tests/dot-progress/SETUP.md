## Preconditions
- The `build` package is importable (`github.com/xhd2015/doctest/libdoc/build`).
- This test creates a temporary sub-tree with 2 leaves (fast + slow) and
  runs `build.Test` on it, capturing stdout to verify dot progress timing.
- Backtick characters in embedded Go strings use `\x60` to avoid
  conflicting with the outer markdown code fence.

## Steps
1. Create sub-tree under a temp dir with `a_fast` (no sleep) and `z_slow`
   (`time.Sleep(5s)` in Setup).
2. Redirect `os.Stdout` to a pipe.
3. In a background goroutine, read the pipe byte-by-byte, recording when the
   first `"."` appears.
4. Call `build.Test(subRoot, core.Options{RemoveTemp: true})`.
5. After `build.Test` returns, close the pipe and collect the result.
6. Return `Incremental` (true if first dot appeared within 4s) and
   `DotCount` (number of dots before the summary line) in the Response.

## Context
- The fast leaf finishes quickly; the slow leaf takes ~5s. If dots are
  incremental, the first dot appears within ~1s (when a_fast completes).
  If batched, all dots appear after ~5s (when z_slow completes).

```go
import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"github.com/xhd2015/doctest/libdoc/build"
	"github.com/xhd2015/doctest/libdoc/core"
)

func writeFile(t *testing.T, root string, rel string, content string) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
```
