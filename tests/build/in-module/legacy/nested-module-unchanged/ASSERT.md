---
label: heavy
---

## Expected

- `doctest test` exits 0 with public import only.
- Nested `go.mod` under outside gen-dir contains `module testcase` and `replace`.
- Generated test imports `example.com/app/pkg/greet`.
- No `.doctest_run_*` temp dirs (legacy path does not use internal import scan).

## Exit Code

- Exit code 0.

```go
import (
	"os"
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
	testData, readErr := os.ReadFile(genTest)
	if readErr != nil {
		t.Fatalf("read generated test: %v", readErr)
	}
	if !strings.Contains(string(testData), modPath+"/pkg/greet") {
		t.Fatalf("expected generated test to import pkg/greet, got:\n%s", string(testData))
	}

	assertNoDoctestRunDirs(t, moduleRoot)
}
```