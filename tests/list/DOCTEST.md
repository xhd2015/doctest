# `doctest list` — inventory of doctest roots (in-process CLI)

## Version
0.0.2

**Layer model (classic TDD — feature not implemented yet; expect RED):**

| Layer | Share | Where |
|-------|-------|--------|
| **L2 doctest in-process** | **all** | every leaf — harness `Run` calls `cli.RunWithWriters` against fixture trees under `t.TempDir()`; no product binary, no `testbin` |

Nested root: does **not** inherit workspace binary Run from `tests/DOCTEST.md`.
All leaves are **unlabeled** (fast). No L3 e2e leaves (short CLI path).

Out of scope (MVP): `--json`, `--paths-only`, L1 package inventory, footer-only mode,
color policy thresholds, process-global `NO_COLOR` env leaves (flag-based color only).

# DSN (Domain Specific Notion)

Inventory of doctest tree roots (dirs with `DOCTEST.md`) with per-root leaf stats.

### Participants

- **Harness** — builds fixture trees under `t.TempDir()`, sets `req.Args`, calls
  `cli.RunWithWriters(&stdout, &stderr, args)` (Parallel-safe writers).
- **CLI dispatcher** — routes `list` (and top-level `--help`) to list inventory
  or usage text.
- **Root discovery** — same notion as `path_resolve` / `test ./...`: patterns are
  plain paths or `path/...`; default `./...`; nested `DOCTEST.md` dirs are separate
  roots; `testdata/` skipped.
- **Leaf inventory** — per root, count `ASSERT.md` leaves (ownership skips nested
  roots and `testdata/`); labels from ASSERT frontmatter; L3 = has label `e2e`.
- **Report writers** — stream one body line per root (tabs), then blank + `---` +
  summary totals/labels; optional gray ANSI on meta when `--color`.

### Behaviors

- **Help** — `list --help` / `-h` documents usage, patterns, L2:L3, color flags;
  top-level `--help` lists `list` among commands.
- **Body line** — `<path>\t<leaves>\tL2:L3=a:b (p2%/p3%)\t<labelDist>`; omit percent
  group when leaves==0; labelDist includes `unlabeled=N`.
- **Summary** — always when ≥1 root: blank, `---`, totals, labels; sums match body.
- **Empty selection** — exit 0, stderr `no tests`, empty stdout (soft exit).
- **Errors** — bare `...`, unknown flags, missing path → non-zero + stderr.
- **Color** — default pipe capture has no ANSI; `--color` grays meta; `--no-color`
  plain; both together → error.

### Pipeline sketch

```
# all leaves (L2 in-process)
fixture under t.TempDir()
  -> req.Args = ["list", flags..., absPattern...]
  -> cli.RunWithWriters(&stdout, &stderr, args)
  -> Response{Stdout, Stderr, ExitCode}
```

## Decision Tree

```
tests/list/
├── DOCTEST.md
├── SETUP.md
├── help/                              dispatch & usage
│   ├── list-help/                     list --help
│   ├── top-level-includes-list/       doctest --help lists list
│   └── unknown-flag/                  list --not-a-real-flag
├── discovery/                         root selection & body shape
│   ├── single-plain-root/
│   ├── multi-root-dotdotdot/          base/... multi-root
│   ├── nested-roots/                  parent+nested; parent excludes nested leaves
│   ├── multi-pattern-union/           union, dedupe, sorted
│   ├── empty-selection/               soft no tests
│   ├── bare-dot-dot-dot/              bare ... rejected
│   ├── testdata-skipped/              testdata not root / not leaves
│   └── missing-path/                  fatal missing path
├── inventory/                         L2:L3 + labelDist semantics
│   ├── unlabeled-all-l2/
│   ├── mix-e2e-unlabeled/
│   ├── multi-label-e2e-heavy/         label: e2e, heavy → L3; both counted
│   ├── heavy-only-is-l2/              heavy without e2e → L2
│   └── zero-leaf-root/                leaves=0, no percent group
├── summary/                           selection-wide footer
│   ├── multi-root-totals/             sums match body; blank+---+totals+labels
│   └── single-root-has-summary/       single root still has footer + trailing \n
└── color/                             flags only (no process Setenv)
    ├── default-no-ansi/               pipe writers → no ANSI
    ├── force-color/                   --color → gray SGR on meta
    ├── force-no-color/                --no-color → no ANSI
    └── conflict/                      --color + --no-color → error
```

## Test Index

| Leaf | Args (sketch) | Expected |
|------|---------------|----------|
| `help/list-help` | `list --help` | usage, patterns, L2:L3, `--color`/`--no-color` |
| `help/top-level-includes-list` | `--help` | command list includes `list` |
| `help/unknown-flag` | `list --not-a-real-flag` | non-zero + usage/error on stderr |
| `discovery/single-plain-root` | `list <root>` | 1 body line + summary |
| `discovery/multi-root-dotdotdot` | `list <base>/...` | ≥2 roots sorted; each body + summary |
| `discovery/nested-roots` | `list <base>/...` | both roots; parent leaves exclude nested |
| `discovery/multi-pattern-union` | `list a b a` | union, dedupe, sorted |
| `discovery/empty-selection` | `list <empty>/...` | exit 0, stderr `no tests`, empty stdout |
| `discovery/bare-dot-dot-dot` | `list ...` | non-zero; bare `...` message |
| `discovery/testdata-skipped` | `list <root>` | no testdata root; leaves ignore testdata ASSERT |
| `discovery/missing-path` | `list <missing>` | non-zero; Error on stderr |
| `inventory/unlabeled-all-l2` | `list <root>` | L2:L3=N:0 (100.0%/0.0%); unlabeled=N |
| `inventory/mix-e2e-unlabeled` | `list <root>` | L3 = e2e count; percents |
| `inventory/multi-label-e2e-heavy` | `list <root>` | L3; e2e=1 heavy=1 |
| `inventory/heavy-only-is-l2` | `list <root>` | L2; heavy in dist |
| `inventory/zero-leaf-root` | `list <root>` | leaves=0; L2:L3=0:0; no percent |
| `summary/multi-root-totals` | `list a b` | totals/labels = sum of body |
| `summary/single-root-has-summary` | `list <root>` | footer present; trailing newline |
| `color/default-no-ansi` | `list <root>` | no ESC in stdout |
| `color/force-color` | `list --color <root>` | gray SGR on meta fields |
| `color/force-no-color` | `list --no-color <root>` | no ESC |
| `color/conflict` | `list --color --no-color` | non-zero conflict error |

## How to Run

```sh
doctest vet ./tests/list/
doctest test ./tests/list/
doctest test ./tests/list/...
```

Expect classic **RED** until `doctest list` is implemented (unknown command / missing inventory).

```go
import (
	"bytes"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/cli"
)

// Request drives one list scenario. Leaves set Args (and optional FixtureDir).
type Request struct {
	// Args is the full argv after the program name, e.g. ["list", root] or ["--help"].
	Args []string

	// FixtureDir is the temp base used by Setup helpers (for asserts that need paths).
	FixtureDir string

	// Roots lists absolute root paths the leaf expects to appear (order not required).
	Roots []string
}

type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Err      error
}

// Run dispatches list/help in-process via cli.RunWithWriters (split stdout/stderr).
// Soft "no tests" and CLI errors are captured into Response — never fail the harness.
// No testbin, no process Setenv/Chdir.
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	var stdout, stderr bytes.Buffer
	args := append([]string(nil), req.Args...)
	err := cli.RunWithWriters(&stdout, &stderr, args)
	resp := &Response{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		Err:    err,
	}
	if err != nil {
		resp.ExitCode = 1
		// Mirror process main: errors not already written to stderr are printed there.
		if strings.TrimSpace(resp.Stderr) == "" {
			resp.Stderr = err.Error() + "\n"
		}
		return resp, nil
	}
	return resp, nil
}
```
