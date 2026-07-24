# `--changed` — policy selection (L2) + sparse process smokes (L3)

## Version
0.0.3

**Layer model (coverage backfill):**

| Layer | Share | Where |
|-------|-------|--------|
| **L2 doctest in-process** | **mass** | `git-context/in-git-repo/**` — fixture trees + synthetic `changedFiles`; harness calls `core.FilterByChangedFiles`, `ChangedRunInfoForTree`, `ChangedDoctestMarkdownFiles` (no product binary, no real git for pure filter policy) |
| **L2 in-process CLI** | **help** | `help/*` — `cli.RunWithWriter` for `<subcmd> --help` documents `--changed`; unlabeled, no product binary |
| **L3 doctest e2e** | **sparse** | `git-context/not-git-repo/*` — real `doctest` binary; `label: heavy` |

Default discovery runs L2 (policy + help). Use `--label heavy` for binary not-git smokes.

Out of scope: product feature changes; `tests/vet`, `tests/help`, `tests/skill`.

# DSN (Domain Specific Notion)

### Participants

- **Harness (L2 policy)** — discovers fixture leaves via `core.DiscoverTreeCases`, applies
  synthetic changed-path lists through `core.FilterByChangedFiles` /
  `ChangedRunInfoForTree` (test/build selection) or
  `ChangedDoctestMarkdownFiles` (vet selection).
- **Harness (L2 help)** — `cli.RunWithWriter` captures subcommand help text that
  documents `--changed` (Parallel-safe; no product binary).
- **e2e binary (L3)** — sparse leaves spawn the product `doctest` binary for
  not-git-repo process errors.
- **Fixture tree** — temporary directory of `DOCTEST.md` / `SETUP.md` /
  `ASSERT.md` (and optional `testdata/`) under `t.TempDir()`; paths are not
  required to live in a real git repo for L2 policy.
- **Changed-file mapper** — maps each changed path to affected leaves: root
  `DOCTEST.md` affects all leaves; group `SETUP.md` affects descendant leaves;
  leaf `ASSERT.md` or `testdata/**` under a leaf affects that leaf only.
  Unrelated non-doctest files under a sibling leaf are ignored.
- **Git gate (L3)** — CLI `--changed` requires a git repository; outside one it
  hard-errors.

### Behaviors

- **Leaf ASSERT change** — only that leaf is selected.
- **Parent/group SETUP change** — all descendant leaves selected.
- **Root DOCTEST change** — every leaf in the tree selected.
- **New untracked leaf** — only the new leaf is selected.
- **Testdata change** — owning leaf selected.
- **Sibling stray file** — does not widen selection beyond truly affected leaves.
- **No matching changes** — zero selected; silent without verbose announce.
- **Build vs test** — share `FilterByChangedFiles` selection policy.
- **Vet** — selects changed doctest markdown only (`ChangedDoctestMarkdownFiles`);
  unchanged invalid root is not selected.
- **Help** — L2 in-process CLI: test/build/vet `--help` documents `--changed`.
- **Not-git** — L3 e2e only.

### Pipeline sketch

```
# L2 (default) — policy
fixture tree under t.TempDir()
  -> core.DiscoverTreeCases(treeDir)
  -> FilterByChangedFiles | ChangedRunInfoForTree | ChangedDoctestMarkdownFiles
       (gitRoot + synthetic changedFiles[] relative to gitRoot)
  -> Response{FilteredPaths, Info, MarkdownPaths}

# L2 (help/) — in-process CLI
req.Args = [<subcmd>, "--help"]
  -> cli.RunWithWriter(&buf, args)
  -> Response{Stdout, ExitCode}

# L3 (not-git-repo/, label: heavy)
testbin.Ensure -> req.Bin
  -> doctest <subcmd> --changed <dir>
```

## Decision Tree

```
tests/changed/
├── DOCTEST.md
├── SETUP.md
├── help/                                      [L2 cli.RunWithWriter, unlabeled]
│   ├── test-help/
│   ├── build-help/
│   └── vet-help/
└── git-context/
    ├── not-git-repo/                          [L3 binary, label: heavy]
    │   ├── test/
    │   ├── build/
    │   └── vet/
    └── in-git-repo/                           [L2 policy, unlabeled]
        ├── test/
        │   ├── assert-only
        │   ├── assert-only-dotdotdot          (same filter; path argv orthogonal)
        │   ├── assert-only-subpath-dotdotdot
        │   ├── sibling-stray-untracked
        │   ├── sibling-stray-subpath-dotdotdot
        │   ├── parent-setup
        │   ├── root-doctest
        │   ├── new-untracked-leaf
        │   ├── testdata-change
        │   └── no-matching-changes/
        │       ├── clean-tree
        │       └── non-doctest-only
        ├── build/
        │   └── assert-only                    (shares FilterByChangedFiles)
        └── vet/
            ├── changed-only                   (ChangedDoctestMarkdownFiles)
            └── skip-root
```

## Test Index

| Leaf | Layer | Expected |
|------|--------|----------|
| `help/test-help` | L2 | stdout includes `Usage: doctest test` and `--changed` |
| `help/build-help` | L2 | stdout includes `Usage: doctest build` and `--changed` |
| `help/vet-help` | L2 | stdout includes `Usage: doctest vet` and `--changed` |
| `not-git-repo/test` | L3 heavy | non-zero; stderr mentions git |
| `not-git-repo/build` | L3 heavy | non-zero; stderr mentions git |
| `not-git-repo/vet` | L3 heavy | non-zero; stderr mentions git |
| `test/assert-only` | L2 | filtered `[leaf_a]`; detail `1 leaf` |
| `test/assert-only-dotdotdot` | L2 | same selection as assert-only |
| `test/assert-only-subpath-dotdotdot` | L2 | same selection as assert-only |
| `test/sibling-stray-untracked` | L2 | filtered `[leaf_a]` despite `leaf_b/stray.go` |
| `test/sibling-stray-subpath-dotdotdot` | L2 | same as sibling-stray-untracked |
| `test/parent-setup` | L2 | both `shared/*` leaves; group SETUP detail |
| `test/root-doctest` | L2 | both leaves; DOCTEST.md detail |
| `test/new-untracked-leaf` | L2 | filtered `[leaf_c]` only |
| `test/testdata-change` | L2 | filtered `[leaf_a]` |
| `test/no-matching-changes/clean-tree` | L2 | ChangedCount 0; silent announce |
| `test/no-matching-changes/non-doctest-only` | L2 | ChangedCount 0 (README outside tree) |
| `build/assert-only` | L2 | same filter as test assert-only |
| `vet/changed-only` | L2 | markdown list is only `leaf_b/SETUP.md` |
| `vet/skip-root` | L2 | markdown list is only `leaf_a/ASSERT.md` (root omitted) |

## How to Run

```sh
doctest vet ./tests/changed/
# default discovery: L2 policy + help (skips label: heavy not-git)
doctest test ./tests/changed/
# sparse binary smokes
doctest test --label heavy ./tests/changed/...
# full suite
doctest test --label-all ./tests/changed/...
```

```go
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/cli"
	"github.com/xhd2015/doctest/libdoc/core"
)

// Policy selects which core changed API Run exercises for L2.
// Empty / "filter" = DiscoverTreeCases + FilterByChangedFiles + ChangedRunInfoForTree.
// "vet-md" = ChangedDoctestMarkdownFiles only (vet selection).
type Policy string

const (
	PolicyFilter Policy = "filter"
	PolicyVetMD  Policy = "vet-md"
)

// Request drives one --changed scenario.
// Default Mode is L2 in-process policy (no binary).
// Set UseCLI for L3 e2e binary only (not-git-repo).
// Help leaves set Args only (no TreeDir, no UseCLI) → cli.RunWithWriter.
type Request struct {
	// L2 policy inputs
	TreeDir      string   // doctest tree root (absolute)
	GitRoot      string   // synthetic or real repo root; changed paths are relative to this
	ChangedFiles []string // paths relative to GitRoot (slash or OS separators OK)
	Policy       Policy   // filter (default) or vet-md

	// L2 help / L3 CLI
	Args    []string
	Env     []string
	WorkDir string
	Timeout time.Duration
	Bin     string
	UseCLI  bool // L3 binary only (not-git-repo)
}

type Response struct {
	// L2 policy outputs
	FilteredPaths []string // TreeCase.Path values, sorted
	Info          core.ChangedRunInfo
	MarkdownPaths []string // paths relative to TreeDir, slash form, sorted
	Announce      bool     // ShouldAnnounceChangedRun(info, verbose=false)

	// L2 help / L3 / shared
	ExitCode int
	Stdout   string
	Stderr   string
	Err      error
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	if req.UseCLI {
		return runE2E(t, req)
	}
	// Help / short-path CLI: Args set without policy TreeDir.
	if len(req.Args) > 0 && req.TreeDir == "" {
		return runCLIWriter(t, req)
	}
	return runPolicy(t, req)
}

// runCLIWriter dispatches help (and similar short-path CLI) via cli.RunWithWriter.
// No testbin, no product binary.
func runCLIWriter(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	var buf bytes.Buffer
	err := cli.RunWithWriter(&buf, req.Args)
	resp := &Response{
		Stdout: buf.String(),
		Err:    err,
	}
	if err != nil {
		resp.ExitCode = 1
		resp.Stderr = err.Error() + "\n"
		return resp, nil
	}
	return resp, nil
}

func runE2E(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	if req.Bin == "" {
		return nil, fmt.Errorf("UseCLI requires req.Bin")
	}
	timeout := req.Timeout
	if timeout == 0 {
		timeout = 45 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, req.Bin, req.Args...)
	if req.WorkDir != "" {
		cmd.Dir = req.WorkDir
	}
	cmd.Env = append(os.Environ(), req.Env...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	runErr := cmd.Run()

	resp := &Response{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Err:      runErr,
		ExitCode: 0,
	}
	if runErr != nil {
		var exitErr *exec.ExitError
		if errors.As(runErr, &exitErr) {
			resp.ExitCode = exitErr.ExitCode()
			return resp, nil
		}
		if ctx.Err() != nil {
			return resp, ctx.Err()
		}
		return resp, runErr
	}
	return resp, nil
}

func runPolicy(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	if req.TreeDir == "" {
		return nil, fmt.Errorf("L2 policy requires req.TreeDir")
	}
	gitRoot := req.GitRoot
	if gitRoot == "" {
		gitRoot = req.TreeDir
	}
	changed := append([]string(nil), req.ChangedFiles...)

	policy := req.Policy
	if policy == "" {
		policy = PolicyFilter
	}

	resp := &Response{}

	switch policy {
	case PolicyVetMD:
		md := core.ChangedDoctestMarkdownFiles(req.TreeDir, gitRoot, changed)
		resp.MarkdownPaths = relMarkdownPaths(req.TreeDir, md)
		return resp, nil
	case PolicyFilter:
		cases, err := core.DiscoverTreeCases(req.TreeDir)
		if err != nil {
			return nil, err
		}
		filtered := core.FilterByChangedFiles(cases, req.TreeDir, gitRoot, changed)
		info := core.ChangedRunInfoForTree(cases, req.TreeDir, gitRoot, changed)
		paths := make([]string, len(filtered))
		for i, tc := range filtered {
			paths[i] = tc.Path
		}
		sort.Strings(paths)
		resp.FilteredPaths = paths
		resp.Info = info
		resp.Announce = core.ShouldAnnounceChangedRun(info, false)
		return resp, nil
	default:
		return nil, fmt.Errorf("unknown Policy %q", policy)
	}
}

func relMarkdownPaths(treeDir string, absPaths []string) []string {
	out := make([]string, 0, len(absPaths))
	// Match core.canonicalAbsPath so macOS /var -> /private/var does not break Rel.
	absTree := canonPath(treeDir)
	for _, p := range absPaths {
		rel, err := filepath.Rel(absTree, canonPath(p))
		if err != nil {
			rel = p
		}
		out = append(out, filepath.ToSlash(rel))
	}
	sort.Strings(out)
	return out
}

func canonPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	if canon, err := filepath.EvalSymlinks(abs); err == nil {
		return canon
	}
	return abs
}

// compile-time touch so core helpers stay importable from leaf helpers if needed.
var (
	_ = core.FilterByChangedFiles
	_ = core.ChangedRunInfoForTree
	_ = core.ChangedDoctestMarkdownFiles
	_ = core.DiscoverTreeCases
	_ = strings.Contains
	_ = cli.Run
)
```
