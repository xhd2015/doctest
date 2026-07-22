# Label Filter Selection (Discover + Filter)

## Version

0.0.2

**Layer L2 in-process** — nested library tree for discovery/filter selection
policy. Leaves build temp fixture trees and call:

- `core.DiscoverTreeCasesLight`
- `core.FilterCasesByLabel` / `core.PartitionLabeledCases`
- `core.FilterBySubDir` when scoping to a leaf path

**No product binary**, **no `label: heavy`**. CLI stdout formatting and
process-boundary contracts live under parent `cli/**` (L3 heavy).

# DSN (Domain Specific Notion)

### Participants

- **Light discover** — walks a fixture DOCTEST tree, reads ASSERT frontmatter
  labels only (`DiscoverTreeCasesLight`).
- **Label filter** — `FilterCasesByLabel` selects run vs skipped cases from
  `LabelExprs`, `LabelAll`, or discovery skip (`PartitionLabeledCases`).
- **Fixture mod** — five-leaf tree: `fast` (unlabeled), `slow`, `ui`, `both`,
  `heavy` (same shape as CLI fixtures).

### Behaviors

- **Single / AND / OR exprs** — only matching *labeled* leaves run; unlabeled
  never match a non-empty expression.
- **Multi-flag OR** — `LabelExprs: ["slow","heavy"]` ≡ OR across flags.
- **No match** — all cases skipped; each skipped has `Reason: "label filter"`.
- **Explicit leaf path** — light-discover whole tree then `FilterBySubDir` to
  one leaf; same label filter applies (match → run, miss → skip).

## Decision Tree

```
select/                               [L2 in-process]
├── filter-single/                    --label slow → slow + both
├── filter-and/                       slow && ui-automation → both
├── filter-or/                        slow || heavy → slow, both, heavy
├── multi-flag-or/                    LabelExprs [slow, heavy] ≡ OR
├── no-match/                         manual → all skipped
├── skip-reason/                      skipped Reason == "label filter"
└── explicit-leaf/
    ├── filter-match/                 subdir slow + --label slow → run
    └── filter-miss/                  subdir slow + --label heavy → skip
```

## Test Index

| Leaf | Expected run paths | Skipped |
|------|--------------------|---------|
| `filter-single` | slow, both | 3 |
| `filter-and` | both | 4 |
| `filter-or` | slow, both, heavy | 2 |
| `multi-flag-or` | slow, both, heavy | 2 |
| `no-match` | (none) | 5 |
| `skip-reason` | (none) | 5 with Reason label filter |
| `explicit-leaf/filter-match` | slow | 0 |
| `explicit-leaf/filter-miss` | (none) | 1 |

## How to Run

```sh
doctest vet ./tests/test/label-filter/select/
# L2 — always discovered (no heavy labels)
doctest test ./tests/test/label-filter/select/...
```

```go
import (
	"fmt"
	"path/filepath"
	"sort"
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
)

// Request configures an in-process discover+filter pass.
type Request struct {
	TreeRoot     string
	LabelExprs   []string
	LabelAll     bool
	ExplicitLeaf bool
	// SubDir, when set, is an absolute path under TreeRoot used with FilterBySubDir
	// (explicit leaf or grouping scope).
	SubDir string
}

type Response struct {
	RunPaths      []string
	SkippedPaths  []string
	SkippedLabels [][]string
	SkippedReason []string
	Err           string
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	if req.TreeRoot == "" {
		return nil, fmt.Errorf("TreeRoot is required")
	}
	cases, err := core.DiscoverTreeCasesLight(req.TreeRoot)
	if err != nil {
		return &Response{Err: err.Error()}, err
	}
	if req.SubDir != "" {
		cases = core.FilterBySubDir(cases, req.TreeRoot, req.SubDir)
	}
	run, skipped := core.FilterCasesByLabel(cases, core.Options{
		LabelExprs:   append([]string(nil), req.LabelExprs...),
		LabelAll:     req.LabelAll,
		ExplicitLeaf: req.ExplicitLeaf,
	})
	resp := &Response{
		RunPaths:      pathsOf(run),
		SkippedPaths:  skippedPathsOf(skipped),
		SkippedLabels: skippedLabelsOf(skipped),
		SkippedReason: skippedReasonsOf(skipped),
	}
	return resp, nil
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

func skippedLabelsOf(skipped []core.SkippedCase) [][]string {
	// Stable order by path.
	type pair struct {
		path   string
		labels []string
	}
	pairs := make([]pair, len(skipped))
	for i, s := range skipped {
		pairs[i] = pair{path: filepath.ToSlash(s.Path), labels: append([]string(nil), s.Labels...)}
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].path < pairs[j].path })
	out := make([][]string, len(pairs))
	for i, p := range pairs {
		out[i] = p.labels
	}
	return out
}

func skippedReasonsOf(skipped []core.SkippedCase) []string {
	type pair struct {
		path, reason string
	}
	pairs := make([]pair, len(skipped))
	for i, s := range skipped {
		pairs[i] = pair{path: filepath.ToSlash(s.Path), reason: s.Reason}
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].path < pairs[j].path })
	out := make([]string, len(pairs))
	for i, p := range pairs {
		out[i] = p.reason
	}
	return out
}
```
