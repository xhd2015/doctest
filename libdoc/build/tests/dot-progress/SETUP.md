## Preconditions
- The `build` package is importable (`github.com/xhd2015/doctest/libdoc/build`).
- This test creates a stable fixture sub-tree with 2 fast leaves and runs
  `build.Test` on it, capturing stdout to verify dot layout.
- Incremental **timing** (dots while a slow package runs) lives in
  `TestDotProgressIncremental` in `build_engine_test.go`.

## Steps
1. Create sub-tree at a stable temp path with two fast leaves.
2. Capture output via `core.Options{Stdout: &buf}` (no global stdout hijacking).
3. Call `build.Test` with stable `GenDir` and `RemoveTemp: false`.
4. Return `DotCount` (dots before the summary line) in the Response.

## Context
- Stable fixture paths let go test cache repeat `doctest test` invocations.
- Structural checks (dot count and placement) stay reliable even when inner
  packages are cached.

```go
import (
	"os"
	"path/filepath"
	"testing"
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