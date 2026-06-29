# `--label` Filter — Doc-Style Test Tree

## Version

0.0.2

# DSN (Domain Specific Notion)

### Participants

- **`doctest test`** — discovers leaves, applies optional `--label EXPR` filter, runs
  matching labeled leaves, skips unlabeled and non-matching labeled leaves, prints skip summary.
- **Label expression parser** — parses `&&`, `||`, parentheses; combines repeatable `--label` flags with OR.
- **Fixture mod tree** — temp tree with `fast` (unlabeled), `slow`, `ui`, `both`, `heavy` leaves.
- **Git fixture** — ephemeral repo for `--changed` then `--label` ordering.

### Behaviors

- **No `--label`** — unchanged discovery skip / explicit-leaf semantics (covered by `label-skip` tree).
- **With `--label`** — only labeled leaves matching EXPR run; all others skipped with `reason: label filter`.
- **Invalid EXPR** — non-zero exit and parse error on stderr before any leaf runs.
- **Help** — `doctest test --help` documents `--label`.
- **Matcher** — `core.EvalLabelExpr(expr, labels)` returns match bool or parse error (library contract).

## Parameter Ranking

| Rank | Factor | Splits at |
|------|--------|-----------|
| 1 | Test surface | `matcher/` (library) vs `cli/` (subprocess) |
| 2 | CLI outcome | `parse-error/` vs `help/` vs `discovery/` vs `explicit-leaf/` vs `with-changed/` |
| 3 | Expression shape | single, AND, OR, precedence, parentheses, whitespace, multi-flag OR |
| 4 | Selection outcome | partial match, single match, no match, skip reason lines |
| 5 | Invocation | tree root discovery vs explicit leaf path |

## Decision Tree

```
label-filter/
├── matcher/                              SURFACE: EvalLabelExpr (nested DOCTEST.md)
│   ├── single-label/
│   ├── and/
│   ├── or/
│   ├── precedence/
│   ├── parentheses/
│   ├── whitespace/
│   └── invalid-syntax/
├── cli/                                  SURFACE: doctest test subprocess
│   ├── parse-error/
│   │   └── trailing-and/
│   ├── help/
│   │   └── documents-label/
│   ├── discovery/
│   │   ├── filter-single/
│   │   ├── filter-and/
│   │   ├── filter-or/
│   │   ├── no-match/
│   │   ├── multi-flag-or/
│   │   └── skip-reason/
│   ├── explicit-leaf/
│   │   ├── filter-match/
│   │   └── filter-miss/
│   └── with-changed/
│       └── changed-then-label/
```

## Test Index

| # | Leaf | Expected |
|---|------|----------|
| 1 | `matcher/single-label/` | `slow` matches `{slow}`; not `{}` or `{fast}` |
| 2 | `matcher/and/` | `slow && ui` matches `{slow,ui}`; not `{slow}` |
| 3 | `matcher/or/` | `slow \|\| heavy` matches `{slow}` and `{heavy}`; not `{fast}` |
| 4 | `matcher/precedence/` | `a \|\| b && c` ≡ `a \|\| (b && c)` |
| 5 | `matcher/parentheses/` | `(slow \|\| heavy) && ui` only when both constraints hold |
| 6 | `matcher/whitespace/` | trimmed ` slow && ui ` parses and matches |
| 7 | `matcher/invalid-syntax/` | `slow &&` parse error |
| 8 | `cli/parse-error/trailing-and/` | non-zero exit, stderr parse error, no PASS line |
| 9 | `cli/help/documents-label/` | stdout includes `--label` |
| 10 | `cli/discovery/filter-single/` | runs `slow` + `both`; skips others |
| 11 | `cli/discovery/filter-and/` | runs `both` only |
| 12 | `cli/discovery/filter-or/` | runs `slow`, `both`, `heavy` |
| 13 | `cli/discovery/no-match/` | exit 0, all skipped, no PASS line |
| 14 | `cli/discovery/multi-flag-or/` | `--label slow --label heavy` ≡ OR expr |
| 15 | `cli/discovery/skip-reason/` | skipped entries include `reason: label filter` |
| 16 | `cli/explicit-leaf/filter-match/` | explicit `slow` + `--label slow` runs |
| 17 | `cli/explicit-leaf/filter-miss/` | explicit `slow` + `--label heavy` skipped |
| 18 | `cli/with-changed/changed-then-label/` | `--changed` then label filter on changed subset |

Regression for behavior **without** `--label` lives in `tests/test/label-skip/` (run separately after implementation).

## How to Run

```sh
doctest vet ./tests/test/label-filter
doctest test ./tests/test/label-filter
doctest test ./tests/test/label-filter/matcher
# regression without --label:
doctest test ./tests/test/label-skip
```

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
	Env     []string
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
	if len(req.Env) > 0 {
		cmd.Env = append(os.Environ(), req.Env...)
	}

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