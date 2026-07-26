# `doctest skill` — list and show (in-process CLI)

## Version
0.0.2

**Layer model (coverage backfill):**

| Layer | Share | Where |
|-------|-------|--------|
| **In-process CLI** | **all** | every leaf — harness `Run` calls `cli.RunWithWriter` (stdout capture) then `cli.Run`; no product binary, no `testbin` |

Nested root: does **not** inherit workspace binary Run from `tests/DOCTEST.md`.
All leaves are **unlabeled** (fast). Completeness: same seven skill scenarios as before.

Out of scope: product feature changes; `tests/vet`, `tests/changed`, `tests/help`,
`skills-update`, skill-show under other trees.

# DSN (Domain Specific Notion)

### Participants

- **Harness** — invokes `cli.RunWithWriter(w, args)` so skill list/show text is
  captured into a buffer (Parallel-safe; does not reassign `os.Stdout`).
- **CLI skill dispatcher (`cli.Run` → `runSkill`)** — lists registered skill
  names or prints embedded skill/spec markdown via `cliStdout`.
- **Embedded skill docs** — TDD, TDD-lite, implementer, doc-spec, code-spec,
  review-perf, and the full skill catalog for `--list`.

### Behaviors

- **List** — `doctest skill --list` prints registered skill names (`doc-spec`,
  `code-spec`, `tdd`, `tdd-cli-agent`, `tdd-lite`, `reproduce`, `review`,
  `review-perf`, `output-assert`, `implementer`, …).
- **Show** — `doctest skill <name> --show` prints the embedded document body
  for that skill (markers differ per leaf).

### Pipeline sketch

```
# all leaves (in-process)
req.Args (e.g. ["skill","--list"] | ["skill","tdd","--show"] | ...)
  -> cli.RunWithWriter(&buf, args)
       -> withTestStdout(buf, cli.Run)
  -> Response{Stdout: buf.String(), ExitCode from err}
```

## Decision Tree

```
tests/skill/
├── DOCTEST.md
├── SETUP.md
├── list/                   skill --list
├── tdd-show/               skill tdd --show
├── tdd-lite-show/          skill tdd-lite --show
├── designer-show/          skill designer --show
├── implementer-show/       skill implementer --show
├── doc-spec-show/          skill doc-spec --show
├── code-spec-show/         skill code-spec --show
└── review-perf-show/       skill review-perf --show
```

## Test Index

| Leaf | Args | Expected markers (subset) |
|------|------|---------------------------|
| `list` | `skill --list` | `doc-spec`, `code-spec`, `tdd`, `tdd-cli-agent`, `tdd-lite`, `reproduce`, `review`, `review-perf`, `output-assert`, `implementer` |
| `tdd-show` | `skill tdd --show` | `doctest-tdd`, `adversarial multi-agent TDD`, plan phases |
| `tdd-lite-show` | `skill tdd-lite --show` | `doctest-tdd-lite`, single-agent cues; no multi-agent orchestrator phrases |
| `designer-show` | `skill designer --show` | `Designer`, `Questions`; no `report-progress` / `yield-pending-questions` |
| `implementer-show` | `skill implementer --show` | `Implementer`, `Questions`; no `report-progress` / `yield-pending-questions` |
| `doc-spec-show` | `skill doc-spec --show` | `doc-style-test-specification`, Directory Layout |
| `code-spec-show` | `skill code-spec --show` | `doc-style-test-code-specification`, Setup/Run/Assert |
| `review-perf-show` | `skill review-perf --show` | `doctest-review-perf`, budgets, metrics flags, WARNING |

## How to Run

```sh
doctest vet ./tests/skill/
doctest test ./tests/skill/
doctest test ./tests/skill/...
```

```go
import (
	"bytes"
	"testing"

	"github.com/xhd2015/doctest/libdoc/cli"
)

// Request drives one skill scenario. Leaves set Args only.
type Request struct {
	Args []string // e.g. ["skill", "--list"], ["skill", "tdd", "--show"]
}

type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Err      error
}

// Run dispatches skill list/show in-process via cli.RunWithWriter (captures cliStdout).
// No testbin, no exec of the product binary.
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
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
		// Mirror process exit: non-zero CLI error is captured, not a harness fail.
		return resp, nil
	}
	return resp, nil
}
```
