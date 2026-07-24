# Short Display Paths — `build.Test` stderr integration

## Version
0.0.2

Integration doc tests verifying that `build.Test` stderr uses `pathfmt.Short`
(display-only) at the real call sites: `→` gen-root announcement, `doctest:`
header, and `cd` command preview — **without process `os.Chdir` / `t.Chdir`**
so both leaves are **Parallel-safe** (`t.Parallel`, `-count=1` full tree).

## DSN (Domain Specific Notion)

### Participants

- **`build.Test`** — discovers doctest leaves, generates Go tests under a gen
  root, prints progress to stderr.
- **`announceRoots`** — prints `→ <genRoot>` as the first stderr line.
- **Progress header** — prints `doctest: <dir>` and test count.
- **`cd` preview** — prints `cd <runDir> && go test ...` before executing.
- **`pathfmt.Short` (DisplayPath)** — display-only formatter at each stderr path
  print: relative to **process cwd** when under cwd, else `~/...` under home,
  else absolute. **Not** Parallel-safe to fake via process Chdir.
- **Harness sandbox** — per-leaf `t.TempDir()` project + absolute `GenDir`
  under that project when explicit; never mutates process cwd.

### Behaviors

- **No process Chdir** — `Run` must not call `os.Chdir` / `t.Chdir`; process cwd
  before and after `build.Test` is unchanged (Parallel-safe lock).
- **Auto gen dir** — mapping-gen cache under home; `cd` line uses `~/...` and
  `mapping-gen` (home shortening does not need project cwd).
- **Explicit gen dir under project** — harness resolves `GenDir` to an
  **absolute** path under the sandbox project (`projRoot/_gen`), not via Chdir
  + relative `"_gen"`. Display is `pathfmt.Short(absGen)` (often absolute under
  temp); still contains the `_gen` segment.
- **Test dir display** — `doctest:` line is `pathfmt.Short(absTestRoot)` (may be
  absolute when sandbox is outside process cwd); still names `tests/feature`.

## Decision Tree

```
short-display-paths
└── gen-dir-source                 [how gen root is chosen]
    ├── mapping-gen-cache          auto cache dir → ~/.../mapping-gen/...
    └── explicit-gen-dir-under-cwd absolute projRoot/_gen (no Chdir)
```

## Test Index

| Leaf | Description | Parallel / Chdir |
|------|-------------|------------------|
| `gen-dir-source/mapping-gen-cache` | Auto gen: `cd` uses `~/` + `mapping-gen`, no raw home; header = `pathfmt.Short(testRoot)`; cwd unchanged | GREEN after harness (no product change) |
| `gen-dir-source/explicit-gen-dir-under-cwd` | Explicit abs `_gen` under sandbox: `cd` uses Short(gen) with `_gen` segment; header Short(testRoot); cwd unchanged | GREEN after harness |

## How to Run

```sh
doctest vet ./libdoc/build/tests/short-display-paths
doctest test -count=1 --label-all ./libdoc/build/tests/short-display-paths
# Both leaves must pass under concurrent t.Parallel (full tree, -count=1).
```

Classic TDD / P2: product non-test code already has no `os.Chdir`. This tree
was the known harness Chdir site. Designer removes Chdir from `Run` and locks
DisplayPath expectations that remain valid without mutating process cwd.

```go
import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/build"
	"github.com/xhd2015/doctest/libdoc/core"
)

type Request struct {
	// GenDir is empty for auto mapping-gen, or a path relative to the
	// per-leaf sandbox project root (e.g. "_gen"). Run joins relative
	// values under projRoot — never relies on process cwd.
	GenDir string
}

type Response struct {
	Stderr         string
	ArrowLine      string
	HeaderLine     string
	CdLine         string
	TestErr        error
	ProjRoot       string
	TestRoot       string
	ResolvedGenDir string // absolute gen dir when explicit; empty when auto
	CwdBefore      string
	CwdAfter       string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	// Parallel-safe: never os.Chdir / t.Chdir. pathfmt.Short uses process
	// Getwd for cwd-relative display; home (~) shortening still works for
	// mapping-gen. Explicit gen dirs are absolute under the sandbox project.
	cwdBefore, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	projRoot := t.TempDir()
	testRoot := createMinimalTree(t, projRoot)

	genDir := req.GenDir
	if genDir != "" && !filepath.IsAbs(genDir) {
		genDir = filepath.Join(projRoot, genDir)
	}

	var stderr bytes.Buffer
	opts := core.Options{
		GenDir:      genDir,
		RemoveTemp:  true,
		Stderr:      &stderr,
	}
	testErr := build.Test(testRoot, opts)
	out := stderr.String()
	arrow, header, cd := parseStderrLines(out)

	cwdAfter, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	return &Response{
		Stderr:         out,
		ArrowLine:      arrow,
		HeaderLine:     header,
		CdLine:         cd,
		TestErr:        testErr,
		ProjRoot:       projRoot,
		TestRoot:       testRoot,
		ResolvedGenDir: genDir,
		CwdBefore:      cwdBefore,
		CwdAfter:       cwdAfter,
	}, nil
}
```
