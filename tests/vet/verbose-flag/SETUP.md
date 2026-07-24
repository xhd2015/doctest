# Scenario

**Feature**: `doctest vet -v` prints directory- and file-level progress (in-process)

```
# injected opts.Stdout (validate.RunWithOptions)
vet -v <dir> -> [vet] validating and SETUP.md lines
```

## Preconditions

- A minimal valid doctest tree is available.
- Verbose progress is captured via injected `opts.Stdout` (no process binary).

## Steps

1. Create a minimal valid doctest tree in a temp directory.
2. Run `vet -v <dir>` via in-process harness (validate + Stdout buffer).

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/libdoc/testtree"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte(testtree.VetDOCTEST()), 0644); err != nil {
		t.Fatal(err)
	}
	// Write SETUP without embedding triple-backtick fences in this go block.
	setupBody := "# Scenario\n\n**Feature**: minimal test setup\n\n" +
		"\x60\x60\x60\n# minimal pipeline\nsystem -> run\n\x60\x60\x60\n\n## Setup\n"
	if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(setupBody), 0644); err != nil {
		t.Fatal(err)
	}
	req.Args = []string{"vet", "-v", dir}
	return nil
}
```
