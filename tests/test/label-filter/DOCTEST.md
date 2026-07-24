# `--label` Filter — Doc-Style Test Tree

## Version

0.0.2

## Layer model (L2 in-process vs L3 e2e)

Coverage-backfill layer split. **Default discovery** runs **L2 only** (unlabeled,
fast library mass). Nested product CLI leaves are **`label: heavy`** and opt-in
via `--label heavy` (or `--label-all`).

| Layer | Subtrees | Run model | Labels |
|-------|----------|-----------|--------|
| **L2 in-process** | `matcher/` (nested), `select/` (nested) | Call `core.EvalLabelExpr`, `DiscoverTreeCasesLight`, `FilterCasesByLabel` — **no** `testbin` | unlabeled |
| **L3 e2e** | `cli/**` | Nested product binary (`testbin.Ensure` + `doctest test`) | **`label: heavy`** + explanation “CLI filter contract via doctest binary” |

| L2 surface | API | Branch |
|------------|-----|--------|
| Expression evaluator | `EvalLabelExpr` | `matcher/**` |
| Discover + selection | `DiscoverTreeCasesLight` + `FilterCasesByLabel` (+ `FilterBySubDir`) | `select/**` |

| L3 surface (sparse) | Contract | Branch |
|---------------------|----------|--------|
| Help text | `doctest test --help` documents `--label` | `cli/help/**` |
| Parse error CLI | invalid EXPR → non-zero, stderr parse/syntax | `cli/parse-error/**` |
| Discovery e2e smoke | full run + compact skip stdout | `cli/discovery/filter-single/` |
| Explicit leaf e2e | path + `--label` runs | `cli/explicit-leaf/filter-match/` |
| Changed + label | `--changed` then label filter | `cli/with-changed/**` |

Pure expression semantics and selection policy live in L2. CLI keeps only
process-boundary contracts (help text, full stdout format, --changed interaction).

# DSN (Domain Specific Notion)

### Participants

- **`doctest test`** — discovers leaves, applies optional `--label EXPR` filter, runs
  matching labeled leaves, skips unlabeled and non-matching labeled leaves, prints skip summary.
- **Label expression parser** — parses `!`, `&&`, `||`, parentheses (precedence `!` > `&&` > `||`); combines repeatable `--label` flags with OR.
- **Fixture mod tree** — temp tree with `fast` (unlabeled), `slow`, `ui`, `both`, `heavy` leaves.
- **Git fixture** — ephemeral repo for `--changed` then `--label` ordering.
- **core selection APIs** — in-process discover/filter used by L2 `select/`.

### Behaviors

- **No `--label`** — unchanged discovery skip / explicit-leaf semantics (covered by `label-skip` tree).
- **With `--label`** — leaves whose label set matches EXPR run (boolean; unlabeled = empty set, so `!e2e` includes them); others skipped in a compact label-filter summary (paths with `-v`).
- **Invalid EXPR** — non-zero exit and parse error on stderr before any leaf runs.
- **Help** — `doctest test --help` documents `--label`.
- **Matcher** — `core.EvalLabelExpr(expr, labels)` returns match bool or parse error (library contract).
- **Select** — light discover + `FilterCasesByLabel` yields run/skipped path sets (library contract).

## Parameter Ranking

| Rank | Factor | Splits at |
|------|--------|-----------|
| 1 | Layer / surface | `matcher/` + `select/` (L2) vs `cli/` (L3) |
| 2 | L2 contract | expression shape vs discover+filter selection |
| 3 | CLI outcome | `parse-error/` vs `help/` vs `discovery/` vs `explicit-leaf/` vs `with-changed/` |
| 4 | Expression shape | single, AND, OR, multi-flag OR, no-match |
| 5 | Invocation | tree root discovery vs explicit leaf path |

## Decision Tree

```
label-filter/
├── matcher/                              [L2] EvalLabelExpr (nested DOCTEST.md)
│   ├── single-label/
│   ├── and/
│   ├── or/
│   ├── precedence/
│   ├── parentheses/
│   ├── whitespace/
│   ├── invalid-syntax/
│   ├── not-bang/
│   ├── not-and/
│   ├── not-parens/
│   └── not-invalid/
├── select/                               [L2] DiscoverTreeCasesLight + FilterCasesByLabel
│   ├── filter-single/
│   ├── filter-and/
│   ├── filter-or/
│   ├── multi-flag-or/
│   ├── no-match/
│   ├── skip-reason/
│   └── explicit-leaf/
│       ├── filter-match/
│       └── filter-miss/
└── cli/                                  [L3 heavy] doctest binary subprocess
    ├── parse-error/
    │   └── trailing-and/
    ├── help/
    │   └── documents-label/
    ├── discovery/
    │   └── filter-single/                e2e smoke: PASS + compact skip
    ├── explicit-leaf/
    │   └── filter-match/
    └── with-changed/
        └── changed-then-label/
```

## Test Index

| # | Leaf | Layer | Expected |
|---|------|-------|----------|
| 1 | `matcher/single-label/` | L2 | `slow` matches `{slow}`; not `{}` or `{fast}` |
| 2 | `matcher/and/` | L2 | `slow && ui` matches `{slow,ui}`; not `{slow}` |
| 3 | `matcher/or/` | L2 | `slow \|\| heavy` matches `{slow}` and `{heavy}`; not `{fast}` |
| 4 | `matcher/precedence/` | L2 | `a \|\| b && c` ≡ `a \|\| (b && c)` |
| 5 | `matcher/parentheses/` | L2 | `(slow \|\| heavy) && ui` only when both constraints hold |
| 6 | `matcher/whitespace/` | L2 | trimmed ` slow && ui ` parses and matches |
| 7 | `matcher/invalid-syntax/` | L2 | `slow &&` parse error |
| 7a | `matcher/not-bang/` | L2 | `!e2e` true on `{}`/`{heavy}`; false on `{e2e}` |
| 7b | `matcher/not-and/` | L2 | `!e2e && heavy` only `{heavy}` |
| 7c | `matcher/not-parens/` | L2 | `!(e2e \|\| flaky)` |
| 7d | `matcher/not-invalid/` | L2 | bare `!` / trailing `!` / `not e2e` error |
| 8 | `select/filter-single/` | L2 | run slow+both |
| 9 | `select/filter-and/` | L2 | run both only |
| 10 | `select/filter-or/` | L2 | run slow, both, heavy |
| 11 | `select/multi-flag-or/` | L2 | LabelExprs OR ≡ single OR expr |
| 12 | `select/no-match/` | L2 | empty run, 5 skips |
| 13 | `select/skip-reason/` | L2 | Reason == "label filter" |
| 14 | `select/explicit-leaf/filter-match/` | L2 | subdir slow matches |
| 15 | `select/explicit-leaf/filter-miss/` | L2 | subdir slow miss → skip |
| 16 | `cli/parse-error/trailing-and/` | L3 | non-zero exit, stderr parse error |
| 17 | `cli/help/documents-label/` | L3 | stdout includes `--label` |
| 18 | `cli/discovery/filter-single/` | L3 | PASS(2/2) + compact skip buckets |
| 19 | `cli/explicit-leaf/filter-match/` | L3 | explicit slow + `--label slow` runs |
| 20 | `cli/with-changed/changed-then-label/` | L3 | `--changed` then label on changed subset |

Regression for behavior **without** `--label` lives in `tests/test/label-skip/` (run separately).

## How to Run

```sh
doctest vet ./tests/test/label-filter
# default discovery: L2 matcher + select (skips label: heavy CLI)
doctest test ./tests/test/label-filter
doctest test ./tests/test/label-filter/matcher/...
doctest test ./tests/test/label-filter/select/...
# L3 CLI e2e
doctest test --label heavy ./tests/test/label-filter/...
# full suite
doctest test --label-all ./tests/test/label-filter/...
# regression without --label:
doctest test ./tests/test/label-skip
```

```go
import (
	"github.com/xhd2015/doctest/libdoc/cli"
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
	UseCLI	bool
	Env     []string
}

type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Err      error
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	if !req.UseCLI {
		var stdout, stderr bytes.Buffer
		err := cli.RunWithWriters(&stdout, &stderr, req.Args)
		resp := &Response{Stdout: stdout.String(), Stderr: stderr.String(), Err: err}
		if err != nil {
			resp.ExitCode = 1
			if resp.Stderr == "" {
				resp.Stderr = err.Error()
			}
			return resp, nil
		}
		return resp, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
	if req.Timeout <= 0 {
		cancel()
		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Minute)
	}
	defer cancel()
	bin := req.Bin
	if bin == "" {
		return nil, fmt.Errorf("UseCLI requires req.Bin")
	}
	cmd := exec.CommandContext(ctx, bin, req.Args...)
	if req.WorkDir != "" {
		cmd.Dir = req.WorkDir
	}
	if len(req.Env) > 0 {
		cmd.Env = append(os.Environ(), req.Env...)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	resp := &Response{Stdout: stdout.String(), Stderr: stderr.String(), Err: err}
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
