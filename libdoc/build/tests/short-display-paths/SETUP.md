# Scenario

**Feature**: `build.Test` stderr paths are shortened via `pathfmt.Short` without process Chdir

```
# build.Test pipeline (Parallel-safe harness)
sandbox projRoot (t.TempDir) + abs GenDir under proj when explicit
  -> build.Test(testRoot, opts)  # no os.Chdir / t.Chdir
  -> announceRoots / doctest header / cd preview
  -> pathfmt.Short(path) on stderr (cwd-rel | ~/home | absolute)

# anti-pattern (forbidden)
Run -> os.Chdir(projRoot)  # process-global; races under t.Parallel
```

## Preconditions

- The `build` package is importable (`github.com/xhd2015/doctest/libdoc/build`).
- Each leaf creates a minimal temp doctest tree and captures stderr from `build.Test`.
- Process cwd is **not** the sandbox project; DisplayPath expectations use
  `pathfmt.Short` on absolute sandbox paths (and home shortening for mapping-gen).
- Backtick characters in embedded Go strings use `\x60` to avoid conflicting with
  the outer markdown code fence.

## Steps

1. Create a temp Go project with `go.mod` and a minimal doctest tree.
2. Resolve relative `req.GenDir` under `projRoot` to an absolute path (no Chdir).
3. Call `build.Test` with stderr captured in a buffer.
4. Parse `→`, `doctest:`, and `cd` lines into `Response`; record cwd before/after.

## Context

- Leaves differ only in `GenDir` / gen-dir source configuration.
- `RemoveTemp: true` avoids leaving generated trees on disk.
- P2 of parallel-safe plan: Chdir anti-pattern (P1 Setenv done; P3 genDir package vars later).

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/testtree"
)

// Helpers only — Run lives in DOCTEST.md and must not Chdir.
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

// assertNoProcessChdir fails if Run mutated process cwd (Parallel anti-pattern).
func assertNoProcessChdir(t *testing.T, resp *Response) {
	t.Helper()
	if resp.CwdBefore == "" || resp.CwdAfter == "" {
		t.Fatal("CwdBefore/CwdAfter must be recorded by Run")
	}
	if resp.CwdBefore != resp.CwdAfter {
		t.Fatalf("Run must not process-Chdir (Parallel-unsafe); cwd %q -> %q", resp.CwdBefore, resp.CwdAfter)
	}
}
```
