# Scenario

**Feature**: `build.Test` stderr paths are shortened via `DisplayPath`

```
# build.Test pipeline
build.Test(dir, opts) -> announceRoots -> doctest header -> cd preview -> go test

# display-only formatting at stderr call sites
genRoot/runDir/dir -> DisplayPath -> shortened stderr line
```

## Preconditions

- The `build` package is importable (`github.com/xhd2015/doctest/libdoc/build`).
- Each leaf creates a minimal temp doctest tree and captures stderr from `build.Test`.
- Backtick characters in embedded Go strings use `\x60` to avoid conflicting with
  the outer markdown code fence.

## Steps

1. Create a temp Go project with `go.mod` and a minimal doctest tree.
2. Optionally `chdir` to the project root so cwd-relative display applies.
3. Call `build.Test` with stderr captured in a buffer.
4. Parse `→`, `doctest:`, and `cd` lines into `Response`.

## Context

- Leaves differ only in `GenDir` / gen-dir source configuration.
- `RemoveTemp: true` avoids leaving generated trees on disk.

```go
import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"github.com/xhd2015/doctest/libdoc/build"
	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/testtree"
)

// Cwd changes for DisplayPath live in root Run under proclock.Mu (not t.Chdir).
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
func createMinimalTree(t *testing.T, projRoot string) string {
	t.Helper()
	if err := os.WriteFile(filepath.Join(projRoot, "go.mod"), []byte("module shortdisplay\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}
	testRoot := filepath.Join(projRoot, "tests", "feature")
	testtree.WriteMinimalRunnableTree(t, testRoot, []testtree.LeafSpec{{Name: "leaf", Steps: "No setup.", Expected: "Passes."}})
	return testRoot
}
func parseStderrLines(stderr string) (string, string, string) {
	var arrow, header, cd string
	for _, line := range strings.Split(stderr, "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "→ "):
			arrow = trimmed
		case strings.HasPrefix(trimmed, "doctest: "):
			header = trimmed
		case strings.HasPrefix(trimmed, "cd "):
			cd = trimmed
		}
	}
	return arrow, header, cd
}
```
