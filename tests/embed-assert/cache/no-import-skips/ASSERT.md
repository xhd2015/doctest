## Expected

- `doctest test` exits 0.
- The current `RawSourceCacheKeyMD5` cache dir is not newly created by this run.

## Exit Code

- Exit code 0.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("doctest subprocess failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected no-import cache skip test to pass, got exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	assertNoAssertReplaceInGoMod(t, filepath.Join(moduleRoot, "go.mod"))
	cacheDir := expectedAssertCacheDir(t)
	_, statErr := os.Stat(cacheDir)
	existsAfter := statErr == nil
	if !cacheDirExistedBefore && existsAfter {
		t.Fatalf("assert-mod cache created without assert import: %s", cacheDir)
	}
}
```