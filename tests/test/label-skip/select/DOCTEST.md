# Labeled Leaf Skip — Selection Policy (in-process)

## Version

0.0.2

**Layer L2 in-process** — nested library tree for discovery skip / label-all /
explicit-leaf selection and frontmatter parse policy. Leaves use:

- `core.DiscoverTreeCasesLight`
- `core.FilterCasesByLabel` / `core.PartitionLabeledCases`
- `core.FilterBySubDir`
- `core.ParseAssertFrontmatter`

**No product binary**, **no `label: e2e`**. CLI summary formatting, path
patterns (`...`), multi-arg, vet/edit/build command surfaces stay under parent
L3 e2e leaves.

# DSN (Domain Specific Notion)

### Participants

- **Light discover** — ASSERT paths + frontmatter labels from a temp tree.
- **Partition / filter** — discovery skip of labeled leaves; `--label-all` runs
  all; explicit leaf path runs labeled leaves (`ExplicitLeaf`).
- **Frontmatter parse** — explanation-only does not set labels (never skip).

### Behaviors

- **Discovery skip** — labeled leaves go to skipped; unlabeled run.
- **All labeled** — run empty, all skipped.
- **Unlabeled only** — run all, skip none.
- **Explanation-only** — no labels → not skipped under discovery.
- **Label-all** — `LabelAll: true` runs labeled + fast.
- **Grouping SubDir** — `FilterBySubDir` then partition skips labeled child.
- **Explicit leaf** — `ExplicitLeaf: true` (no LabelExprs) runs labeled leaf.
- **Frontmatter** — parse yields labels / explanation; malformed errors.

## Decision Tree

```
select/                                   [L2 in-process]
├── mixed-fast-labeled/                   1 run + 1 skipped
├── all-labeled/                          0 run, all skipped
├── unlabeled-only/                       all run, no skip
├── explanation-only-runs/                explanation w/o label runs
├── label-all-runs-all/                   LabelAll → both run
├── grouping-dir/                         SubDir e2e → skip labeled child
├── explicit-leaf-runs/                   ExplicitLeaf runs labeled
└── frontmatter/
    ├── explanation-only/                 parse: labels empty
    └── labeled/                          parse: labels + explanation
```

## Test Index

| Leaf | Expected |
|------|----------|
| `mixed-fast-labeled` | run {fast_leaf}; skip {labeled_leaf} |
| `all-labeled` | run {}; skip {labeled_leaf} |
| `unlabeled-only` | run {plain_leaf}; skip {} |
| `explanation-only-runs` | run {explained_leaf}; skip {} |
| `label-all-runs-all` | run both |
| `grouping-dir` | run {e2e/fast_child}; skip {e2e/labeled_child} |
| `explicit-leaf-runs` | ExplicitLeaf on labeled → run |
| `frontmatter/explanation-only` | Labels empty; Explanation set |
| `frontmatter/labeled` | Labels non-empty |

## How to Run

```sh
doctest vet ./tests/test/label-skip/select/
doctest test ./tests/test/label-skip/select/...
```

```go
import (
	"fmt"
	"path/filepath"
	"sort"
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
)

// Request selects one in-process surface.
// Op: "filter" (default) | "parse_frontmatter"
type Request struct {
	Op string

	// filter
	TreeRoot     string
	LabelExprs   []string
	LabelAll     bool
	ExplicitLeaf bool
	SubDir       string

	// parse_frontmatter
	FrontmatterContent string
	FrontmatterPath    string
}

type Response struct {
	RunPaths     []string
	SkippedPaths []string

	// parse_frontmatter
	Labels      []string
	Explanation string
	ParseErr    string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	op := req.Op
	if op == "" {
		op = "filter"
	}
	switch op {
	case "filter":
		if req.TreeRoot == "" {
			return nil, fmt.Errorf("TreeRoot is required")
		}
		cases, err := core.DiscoverTreeCasesLight(req.TreeRoot)
		if err != nil {
			return nil, err
		}
		if req.SubDir != "" {
			cases = core.FilterBySubDir(cases, req.TreeRoot, req.SubDir)
		}
		run, skipped := core.FilterCasesByLabel(cases, core.Options{
			LabelExprs:   append([]string(nil), req.LabelExprs...),
			LabelAll:     req.LabelAll,
			ExplicitLeaf: req.ExplicitLeaf,
		})
		return &Response{
			RunPaths:     pathsOf(run),
			SkippedPaths: skippedPathsOf(skipped),
		}, nil

	case "parse_frontmatter":
		path := req.FrontmatterPath
		if path == "" {
			path = "ASSERT.md"
		}
		fm, _, err := core.ParseAssertFrontmatter(path, req.FrontmatterContent)
		resp := &Response{
			Labels:      append([]string(nil), fm.Labels...),
			Explanation: fm.Explanation,
		}
		if err != nil {
			resp.ParseErr = err.Error()
		}
		return resp, nil

	default:
		return nil, fmt.Errorf("unknown Op %q", op)
	}
}

func pathsOf(cases []core.TreeCase) []string {
	out := make([]string, len(cases))
	for i, c := range cases {
		out[i] = filepath.ToSlash(c.Path)
	}
	sort.Strings(out)
	return out
}

func skippedPathsOf(skipped []core.SkippedCase) []string {
	out := make([]string, len(skipped))
	for i, s := range skipped {
		out[i] = filepath.ToSlash(s.Path)
	}
	sort.Strings(out)
	return out
}
```
