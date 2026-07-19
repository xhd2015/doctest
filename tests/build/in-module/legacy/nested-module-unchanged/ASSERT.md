---
label: heavy
---

## Expected

- `doctest test` exits 0 with public import only.
- Nested `go.mod` under outside gen-dir contains `module testcase` and `replace`.
- Generated tree imports `example.com/app/pkg/greet` (in root package under unified layout).
- No `.doctest_run_*` temp dirs (legacy path does not use internal import scan).

## Exit Code

- Exit code 0.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("doctest subprocess failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected legacy nested-module test to pass, got exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}

	assertNestedGoMod(t, outsideGenDir)

	genTest := generatedLeafTestPath(outsideGenDir)
	assertFileExists(t, genTest)

	// Unified: Run (and its imports) live in __droot; leaf is thin.
	wantImport := modPath + "/pkg/greet"
	found := false
	filepath.Walk(outsideGenDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		if strings.Contains(string(data), wantImport) {
			found = true
		}
		return nil
	})
	if !found {
		t.Fatalf("expected generated tree to import %s under %s", wantImport, outsideGenDir)
	}

	assertNoDoctestRunDirs(t, moduleRoot)
}
```