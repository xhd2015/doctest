# Labeled Leaf Skip — Doc-Style Test Tree

## Version

0.0.2

## Layer model (L2 in-process vs L3 e2e)

Coverage-backfill layer split. **Default discovery** runs **L2 only** (unlabeled
`select/**`). Product CLI leaves under `test/`, `vet/`, `edit/`, `build/` are
**`label: heavy`** and opt-in via `--label heavy`.

| Layer | Subtrees | Run model | Labels |
|-------|----------|-----------|--------|
| **L2 in-process** | `select/` (nested) | `DiscoverTreeCasesLight` + `FilterCasesByLabel` / `ParseAssertFrontmatter` — **no** `testbin` | unlabeled |
| **L3 e2e** | `test/**`, `vet/**`, `edit/**`, `build/**` | Product binary subprocess | **`label: heavy`** |

| L2 surface | API | Branch |
|------------|-----|--------|
| Discovery skip / label-all / explicit | `FilterCasesByLabel`, `FilterBySubDir` | `select/*` (except frontmatter/) |
| Frontmatter parse | `ParseAssertFrontmatter` | `select/frontmatter/**` |

| L3 surface (sparse process contracts) | Branch |
|---------------------------------------|--------|
| Skip summary stdout + all-skip no PASS | `test/discovery/mixed-fast-labeled`, `all-labeled` |
| `./...` pattern, multi-arg | `test/discovery/dotdotdot-pattern`, `multi-arg-mixed` |
| `--label-all` vs `--label` mutual exclusion | `test/discovery/label-all-conflicts-label` |
| Explicit leaf / ASSERT.md path run | `test/explicit-leaf/**` |
| `doctest vet` frontmatter validation | `vet/**` |
| `doctest edit` mutations | `edit/**` |
| `doctest build` compiles labeled | `build/**` |

# DSN (Domain Specific Notion)

### Participants

- **`doctest test`** — discovers runnable leaves, skips labeled ones in discovery mode, runs them for explicit leaf paths, prints a skip summary.
- **`doctest build`** — compiles all leaves including labeled ones (skip does not apply).
- **`doctest vet`** — validates ASSERT.md frontmatter YAML on every leaf.
- **`doctest edit`** — updates frontmatter on a single concrete leaf ASSERT.md.
- **Temp test tree** — programmatic fixture with fast and/or labeled leaves.
- **core selection APIs** — in-process discover/filter used by L2 `select/`.

### Behaviors

- **Discovery skip** — tree root, grouping dir, or `...` pattern omits leaves whose ASSERT.md has `label:`.
- **`--label-all`** — discovery runs labeled leaves too (full suite); mutually exclusive with `--label`.
- **Explicit leaf** — concrete leaf directory or `ASSERT.md` path runs labeled tests.
- **Explanation-only** — frontmatter with `explanation:` but no `label:` never skips.
- **Skip summary** — compact by label set (counts); paths/explanations with `-v`.
- **All skipped** — exit 0 when every discovered leaf is labeled.
- **Build all** — `doctest build` compiles labeled leaves even when `doctest test` would skip them.
- **Edit leaf** — add label/explanation; reject `...` patterns.

## Parameter Ranking

| Rank | Factor | Splits at |
|------|--------|-----------|
| 1 | Layer | `select/` (L2) vs CLI commands (L3) |
| 2 | CLI command | `test/`, `vet/`, `edit/`, `build/` |
| 3 | Invocation mode (test only) | `discovery/` vs `explicit-leaf/` |
| 4 | Discovery outcome | mixed, all-labeled, `...` pattern, multi-arg, flag conflicts |
| 5 | Vet / edit / build | validity, mutation, tree shape |

## Decision Tree

```
label-skip/
├── select/                               [L2] core partition/filter + frontmatter parse
│   ├── mixed-fast-labeled/
│   ├── all-labeled/
│   ├── unlabeled-only/
│   ├── explanation-only-runs/
│   ├── label-all-runs-all/
│   ├── grouping-dir/
│   ├── explicit-leaf-runs/
│   └── frontmatter/
│       ├── explanation-only/
│       └── labeled/
├── test/                                 [L3 heavy] doctest test
│   ├── discovery/
│   │   ├── mixed-fast-labeled/           OUTCOME: 1 run + skip summary format
│   │   ├── all-labeled/                  OUTCOME: 0 run, skip summary, no PASS
│   │   ├── dotdotdot-pattern/            OUTCOME: ./mod/... path pattern
│   │   ├── multi-arg-mixed/              OUTCOME: multi-arg aggregation
│   │   └── label-all-conflicts-label/    OUTCOME: mutual exclusion error
│   └── explicit-leaf/
│       ├── runs-labeled/
│       └── assert-md-path/
├── vet/                                  [L3 heavy] doctest vet
│   ├── valid-frontmatter/
│   ├── explanation-only/
│   └── malformed-frontmatter/
├── edit/                                 [L3 heavy] doctest edit
│   ├── add-label/
│   ├── set-label-on-existing-frontmatter/
│   ├── append-explanation/
│   ├── idempotent-label-warn/
│   ├── assert-md-path/
│   └── rejects-dotdotdot/
└── build/                                [L3 heavy] doctest build
    ├── compiles-labeled/
    └── mixed-tree/
```

## Test Index

| # | Leaf | Layer | Expected |
|---|------|-------|----------|
| 1 | `select/mixed-fast-labeled/` | L2 | run fast; skip labeled |
| 2 | `select/all-labeled/` | L2 | all skipped |
| 3 | `select/unlabeled-only/` | L2 | all run |
| 4 | `select/explanation-only-runs/` | L2 | explanation w/o label runs |
| 5 | `select/label-all-runs-all/` | L2 | LabelAll runs both |
| 6 | `select/grouping-dir/` | L2 | SubDir e2e skips labeled child |
| 7 | `select/explicit-leaf-runs/` | L2 | ExplicitLeaf runs labeled |
| 8 | `select/frontmatter/explanation-only/` | L2 | parse: empty labels |
| 9 | `select/frontmatter/labeled/` | L2 | parse: labels set |
| 10 | `test/discovery/mixed-fast-labeled/` | L3 | PASS(1/1) + compact skip block |
| 11 | `test/discovery/all-labeled/` | L3 | exit 0, skip block, no PASS line |
| 12 | `test/discovery/dotdotdot-pattern/` | L3 | `./mod/...` skip format |
| 13 | `test/discovery/multi-arg-mixed/` | L3 | multi-arg PASS + skip |
| 14 | `test/discovery/label-all-conflicts-label/` | L3 | mutual exclusion error |
| 15 | `test/explicit-leaf/runs-labeled/` | L3 | labeled leaf dir executes |
| 16 | `test/explicit-leaf/assert-md-path/` | L3 | ASSERT.md path executes |
| 17 | `vet/valid-frontmatter/` | L3 | exit 0 |
| 18 | `vet/explanation-only/` | L3 | exit 0 |
| 19 | `vet/malformed-frontmatter/` | L3 | non-zero exit |
| 20 | `edit/add-label/` | L3 | exact frontmatter after edit |
| 21 | `edit/set-label-on-existing-frontmatter/` | L3 | comma-separated labels |
| 22 | `edit/append-explanation/` | L3 | `first; second` |
| 23 | `edit/idempotent-label-warn/` | L3 | stderr warning, no change |
| 24 | `edit/assert-md-path/` | L3 | edit via ASSERT.md path |
| 25 | `edit/rejects-dotdotdot/` | L3 | non-zero exit |
| 26 | `build/compiles-labeled/` | L3 | labeled-only tree compiles |
| 27 | `build/mixed-tree/` | L3 | fast + labeled both compile |

## How to Run

```sh
doctest vet ./tests/test/label-skip
# default discovery: L2 select only (skips label: heavy CLI)
doctest test ./tests/test/label-skip
doctest test ./tests/test/label-skip/select/...
# L3 CLI e2e
doctest test --label heavy ./tests/test/label-skip/...
# full suite
doctest test --label-all ./tests/test/label-skip/...
```

```go
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"testing"
	"time"
)

type Request struct {
	Args    []string
	WorkDir string
	Timeout time.Duration
	Bin     string
}

type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Err      error
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
	defer cancel()

	if req.Bin == "" {
		return nil, fmt.Errorf("req.Bin is not set")
	}
	cmd := exec.CommandContext(ctx, req.Bin, req.Args...)
	cmd.Dir = req.WorkDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	resp := &Response{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		Err:    err,
	}
	if err == nil {
		return resp, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		resp.ExitCode = exitErr.ExitCode()
		return resp, nil
	}
	if ctx.Err() != nil {
		return resp, ctx.Err()
	}
	return resp, err
}

```
