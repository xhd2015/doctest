# Short Display Paths — `build.Test` stderr integration

## Version
0.0.2


Integration doc tests verifying that `build.Test` stderr uses `DisplayPath` at
the real call sites: `→` gen-root announcement, `doctest:` header, and `cd`
command preview.

## DSN (Domain Specific Notion)

### Participants

- **`build.Test`** — discovers doctest leaves, generates Go tests under a gen
  root, prints progress to stderr.
- **`announceRoots`** — prints `→ <genRoot>` as the first stderr line.
- **Progress header** — prints `doctest: <dir>` and `─── N test cases`.
- **`cd` preview** — prints `cd <runDir> && go test ...` before executing.
- **`DisplayPath`** — display-only formatter applied at each stderr path print.

### Behaviors

- **Auto gen dir** — mapping-gen cache under home; `→` and `cd` lines use `~/...`
  instead of `/Users/...`.
- **Explicit gen dir under cwd** — user `--gen-dir` under project; `→` line uses
  `_gen` when cwd is the project root.
- **Test dir under cwd** — `doctest:` line uses cwd-relative path without `./` prefix.

## Decision Tree

```
short-display-paths
└── gen-dir-source              [how gen root is chosen]
    ├── mapping-gen-cache       auto cache dir → ~/.../mapping-gen/...
    └── explicit-gen-dir-under-cwd  --gen-dir _gen under project
```

## Test Index

| Leaf | Description |
|------|-------------|
| `gen-dir-source/mapping-gen-cache` | Auto gen dir stderr uses `~/...mapping-gen...`; `doctest:` uses cwd-relative path; no raw home absolute |
| `gen-dir-source/explicit-gen-dir-under-cwd` | Explicit `_gen` under project displays as `→ _gen` |

## How to Run

```sh
doctest vet ./libdoc/build/tests/short-display-paths/...
doctest test ./libdoc/build/tests/short-display-paths/...
```

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
	GenDir string
}
type Response struct {
	Stderr		string
	ArrowLine	string
	HeaderLine	string
	CdLine		string
	TestErr		error
}
func Run(t *testing.T, req *Request) (*Response, error) {
	// DisplayPath shortens relative to process cwd. This tree needs a real
	// project cwd; os.Chdir is process-global (not Parallel-safe across trees).
	// Prefer labeling this tree out of concurrent light-tree Parallel, or
	// redesign DisplayPath tests without cwd. Never t.Chdir.
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(wd) }()

	projRoot := t.TempDir()
	testRoot := createMinimalTree(t, projRoot)
	if err := os.Chdir(projRoot); err != nil {
		t.Fatal(err)
	}

	var stderr bytes.Buffer
	opts := core.Options{
		GenDir:		req.GenDir,
		RemoveTemp:	true,
		Stderr:		&stderr,
	}
	testErr := build.Test(testRoot, opts)
	out := stderr.String()
	arrow, header, cd := parseStderrLines(out)

	return &Response{
		Stderr:		out,
		ArrowLine:	arrow,
		HeaderLine:	header,
		CdLine:		cd,
		TestErr:	testErr,
	}, nil
}
```
