---
label: heavy
---

## Expected

- `doctest test` exits 0.
- Gen-dir dump contains generated `leaf_test.go` importing internal greet and assert.
- Dump has no nested `go.mod`.
- Compile temp removed after run.

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
		t.Fatalf("expected internal+assert gen-dir dump test to pass, got exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	dumpDir := filepath.Join(moduleRoot, "_gen")
	genTest := generatedLeafTestPath(dumpDir)
	assertFileExists(t, genTest)
	data, err := os.ReadFile(genTest)
	if err != nil {
		t.Fatalf("read dump test file: %v", err)
	}
	src := string(data)
	if !strings.Contains(src, modPath+"/internal/greet") {
		t.Fatalf("expected dump to import internal/greet, got:\n%s", src)
	}
	if !strings.Contains(src, assertModPath) {
		t.Fatalf("expected dump to import assert, got:\n%s", src)
	}
	assertDumpNoNestedGoMod(t, dumpDir)
	assertStderrUsesTempCompile(t, resp)
	assertNoDoctestRunDirs(t, moduleRoot)
}
```