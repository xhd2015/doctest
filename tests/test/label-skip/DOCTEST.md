# Labeled Leaf Skip — Doc-Style Test Tree

## Version

0.0.2

# DSN (Domain Specific Notion)

### Participants

- **`doctest test`** — discovers runnable leaves, skips labeled ones in discovery mode, runs them for explicit leaf paths, prints a skip summary.
- **`doctest build`** — compiles all leaves including labeled ones (skip does not apply).
- **`doctest vet`** — validates ASSERT.md frontmatter YAML on every leaf.
- **`doctest edit`** — updates frontmatter on a single concrete leaf ASSERT.md.
- **Temp test tree** — programmatic fixture with fast and/or labeled leaves.

### Behaviors

- **Discovery skip** — tree root, grouping dir, or `...` pattern omits leaves whose ASSERT.md has `label:`.
- **Explicit leaf** — concrete leaf directory or `ASSERT.md` path runs labeled tests.
- **Explanation-only** — frontmatter with `explanation:` but no `label:` never skips.
- **Skip summary** — lists skipped path, label, and explanation before PASS/FAIL.
- **All skipped** — exit 0 when every discovered leaf is labeled.
- **Build all** — `doctest build` compiles labeled leaves even when `doctest test` would skip them.
- **Edit leaf** — add label/explanation; reject `...` patterns.

## Parameter Ranking

| Rank | Factor | Splits at |
|------|--------|-----------|
| 1 | CLI command | `test/`, `vet/`, `edit/`, `build/` |
| 2 | Invocation mode (test only) | `discovery/` vs `explicit-leaf/` |
| 3 | Discovery outcome | mixed, all-labeled, unlabeled-only, explanation-only, grouping dir, `...` pattern, multi-arg |
| 4 | Vet input validity | valid frontmatter, explanation-only, malformed YAML |
| 5 | Edit mutation / input | add label, append label, append explanation, idempotent warn, ASSERT.md path, reject `...` |
| 6 | Build tree shape | labeled-only, mixed fast+labeled |

## Decision Tree

```
label-skip/
├── test/                                 COMMAND: doctest test
│   ├── discovery/                        MODE: tree root, grouping, ./.../, or multi-arg
│   │   ├── mixed-fast-labeled/           OUTCOME: 1 run + 1 skipped
│   │   ├── all-labeled/                  OUTCOME: 0 run, all skipped, exit 0
│   │   ├── unlabeled-only/               OUTCOME: PASS, no skip block
│   │   ├── explanation-only-runs/        OUTCOME: explanation without label runs
│   │   ├── grouping-dir/                 OUTCOME: grouping dir skips labeled child
│   │   ├── dotdotdot-pattern/            OUTCOME: ./mod/... skips labeled
│   │   └── multi-arg-mixed/              OUTCOME: ./mod/... + explicit leaf, aggregated skip
│   └── explicit-leaf/                    MODE: concrete leaf dir or ASSERT.md
│       ├── runs-labeled/                 OUTCOME: labeled leaf dir executes
│       └── assert-md-path/               OUTCOME: labeled leaf via ASSERT.md path
├── vet/                                  COMMAND: doctest vet
│   ├── valid-frontmatter/                INPUT: label+explanation → exit 0
│   ├── explanation-only/                 INPUT: explanation without label → exit 0
│   └── malformed-frontmatter/            INPUT: broken YAML → non-zero exit
├── edit/                                 COMMAND: doctest edit
│   ├── add-label/                        MUTATION: create frontmatter with label
│   ├── set-label-on-existing-frontmatter/ MUTATION: append second label
│   ├── append-explanation/               MUTATION: append with "; "
│   ├── idempotent-label-warn/            MUTATION: duplicate label warns, no change
│   ├── assert-md-path/                   INPUT: ASSERT.md path accepted
│   └── rejects-dotdotdot/                INPUT: ... path → error
└── build/                                COMMAND: doctest build
    ├── compiles-labeled/                 OUTCOME: labeled-only tree compiles
    └── mixed-tree/                       OUTCOME: fast + labeled both compile
```

## Test Index

| # | Leaf | Expected |
|---|------|----------|
| 1 | `test/discovery/mixed-fast-labeled/` | PASS(1/1) + exact skip block for labeled leaf |
| 2 | `test/discovery/all-labeled/` | exit 0, exact skip block, no PASS/FAIL line |
| 3 | `test/discovery/unlabeled-only/` | PASS(1/1), no skip block |
| 4 | `test/discovery/explanation-only-runs/` | PASS(1/1), no skip block |
| 5 | `test/discovery/grouping-dir/` | PASS(1/1) + skip labeled child under grouping dir |
| 6 | `test/discovery/dotdotdot-pattern/` | PASS(1/1) + skip block via `./mod/...` |
| 7 | `test/discovery/multi-arg-mixed/` | PASS(2/2) + aggregated skip from discovery pass |
| 8 | `test/explicit-leaf/runs-labeled/` | PASS(1/1), no skip block |
| 9 | `test/explicit-leaf/assert-md-path/` | PASS(1/1) via ASSERT.md path, no skip block |
| 10 | `vet/valid-frontmatter/` | exit 0 |
| 11 | `vet/explanation-only/` | exit 0 |
| 12 | `vet/malformed-frontmatter/` | non-zero exit |
| 13 | `edit/add-label/` | ASSERT.md exact frontmatter after edit |
| 14 | `edit/set-label-on-existing-frontmatter/` | `label: ui-automation, manual` exact |
| 15 | `edit/append-explanation/` | explanation `first; second` exact |
| 16 | `edit/idempotent-label-warn/` | stderr warning exact, ASSERT.md unchanged |
| 17 | `edit/assert-md-path/` | edit via ASSERT.md path, exact frontmatter |
| 18 | `edit/rejects-dotdotdot/` | non-zero exit |
| 19 | `build/compiles-labeled/` | exit 0, compiles labeled-only tree |
| 20 | `build/mixed-tree/` | exit 0, compiles fast + labeled tree |

## How to Run

Run `doctest vet ./tests/test/label-skip` then `doctest test ./tests/test/label-skip`.

```go
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
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

func Run(t *testing.T, req *Request) (*Response, error) {
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