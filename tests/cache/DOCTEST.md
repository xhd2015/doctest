# `doctest cache [--clean] [--dry-run]` — cache info and clean (L2 in-process)

## Version
0.0.2

**Layer model (classic TDD — feature not implemented yet; expect RED):**

| Layer | Share | Where |
|-------|-------|--------|
| **L2 in-process CLI** | **~80%** | every leaf — harness `Run` calls `cli.RunWithWriters`; help, flag errors, info, clean |
| **L2 injectable roots** | **~20%** | `Request.CacheHome` / override paths from `t.TempDir()`; no process `Setenv`/`Chdir` |
| **L3 e2e** | **avoid** | no product binary / `testbin` |

Nested root: does **not** inherit workspace binary Run from `tests/DOCTEST.md`.
All leaves are **unlabeled** (fast).

Out of scope: `GOCACHE`, go testcache, selective `--only`, JSON, `cache path`
subcommand, process-global env leaves.

# DSN (Domain Specific Notion)

Inventory and wipe of durable doctest cache roots under `$CacheHome/doctest`
(plus override roots that live outside that tree).

### Participants

- **Harness** — isolates a temp `CacheHome` (and optional outside overrides) on
  `Request`, seeds bucket files, sets `Args`, calls `cli.RunWithWriters`
  (Parallel-safe; never reassigns `os.Stdout` / process env).
- **CLI dispatcher** — routes `cache` (and top-level `--help`) to info, clean,
  or usage / flag errors.
- **Cache home** — base directory (`DOCTEST_CACHE_HOME` or user cache dir in
  production; **injected** as `Request.CacheHome` in tests).
- **Doctest root** — `$CacheHome/doctest` (entire tree is the clean target).
- **Buckets** — first-level subdirs under doctest root (`leaf-cache`,
  `mapping-gen`, `metrics`, `assert-mod`, `session-mod`, …).
- **Override roots** — paths outside the main doctest root that product also
  owns when env overrides point there (`DOCTEST_LEAF_CACHE` store root;
  metrics under `$MetricsRoot/doctest/metrics` when `DOCTEST_METRICS_ROOT`
  base is outside CacheHome). Tests inject absolute override paths on
  `Request` (no process Setenv).
- **Human size** — byte counts rendered with K/M/G (and B for tiny totals).

### Behaviors

- **Info (default)** — `doctest cache`: print Cache home, Doctest root, per-bucket
  human sizes, Total + bucket count. Empty / missing root still prints home +
  root and indicates empty (0B / 0 buckets / no cache).
- **Clean live** — `doctest cache --clean`: remove entire doctest root (+ outside
  overrides when set); stdout `Removed <abs>  (<human-size>)` per target.
- **Clean dry-run** — `doctest cache --clean --dry-run`: print
  `[dry-run] would remove: <abs>  (<human-size>)`; **no** deletes.
- **Dry-run alone** — `doctest cache --dry-run`: non-zero; message that
  `--dry-run requires --clean`.
- **Help** — `cache --help` / `-h`: usage on stdout mentioning `--clean` and
  `--dry-run`; exit 0. Top-level `--help` lists `cache`.
- **Unknown flag** — non-zero + error/usage (must be a real `cache` flag error,
  not merely `unknown command: cache` after registration).
- **Safety** — refuse remove if path empty, not absolute after resolve,
  filesystem root, or main root does not end with path component `doctest`.
- **Clean failure** — hard-fail on remove error (non-zero + path); no
  partial-success exit 0.

### Pipeline sketch

```
# all leaves (L2 in-process)
req.CacheHome = t.TempDir()  (+ optional LeafCache / MetricsRoot outside)
  -> seed buckets under CacheHome/doctest/...
  -> req.Args = ["cache"] | ["cache","--clean"] | ["cache","--clean","--dry-run"] | ...
  -> cli.RunWithWriters(&stdout, &stderr, args)
       // implementer: wire CacheHome / overrides without process Setenv
       // (package Scan/Clean opts preferred; CLI production still reads env)
  -> Response{Stdout, Stderr, ExitCode}
  -> Assert: output shape + filesystem side effects under req paths
```

## Decision Tree

```
tests/cache/
├── DOCTEST.md
├── SETUP.md
├── help/                              dispatch & usage
│   ├── cache-help/                    cache --help; mentions --clean, --dry-run
│   └── top-level-includes-cache/      doctest --help lists cache
├── info/                              default cache (no --clean)
│   ├── empty-root/                    empty/missing doctest root
│   └── seeded-buckets/                ≥2 buckets; human sizes + Total
├── clean/                             --clean (+ dry-run matrix)
│   ├── dry-run-requires-clean/        --dry-run alone → error
│   ├── dry-run/                       --clean --dry-run; no delete
│   └── live/                          --clean removes tree
├── flags/                             invalid argv
│   └── unknown-flag/                  cache --not-a-real-flag
└── overrides/                         roots outside main doctest tree
    └── leaf-cache-outside/            outside leaf-cache in dry-run plan
```

Split factor at root: **invocation mode** (help | info | clean | flags |
overrides). Under `clean/`, siblings partition dry-run-only error vs dry-run
preview vs live remove. Under `info/`, empty vs seeded content.

## Test Index

| Leaf | Args (sketch) | Expected |
|------|---------------|----------|
| `help/cache-help` | `cache --help` | exit 0; usage; `--clean`, `--dry-run` |
| `help/top-level-includes-cache` | `--help` | command list includes `cache` |
| `info/empty-root` | `cache` + empty home | exit 0; Cache home + Doctest root; empty/0 |
| `info/seeded-buckets` | `cache` + ≥2 buckets | bucket names; Total non-zero; human units |
| `clean/dry-run-requires-clean` | `cache --dry-run` | non-zero; requires `--clean` |
| `clean/dry-run` | seed + `cache --clean --dry-run` | `[dry-run] would remove`; tree still exists |
| `clean/live` | seed + `cache --clean` | tree gone; stdout `Removed` |
| `flags/unknown-flag` | `cache --not-a-real-flag` | non-zero; flag/usage error (not only unknown command) |
| `overrides/leaf-cache-outside` | outside leaf-cache + clean dry-run | both main root and override in would-remove |

## How to Run

```sh
doctest vet ./tests/cache/
doctest test ./tests/cache/
doctest test ./tests/cache/...
```

Expect classic **RED** until `doctest cache` is implemented (unknown command /
missing behavior).

### How Run works (injection strategy)

1. **CLI surface** — every leaf goes through `cli.RunWithWriters(&stdout, &stderr, req.Args)`
   so help text, flag errors, and command routing match the product binary
   (Parallel-safe writers; no `os.Stdout` reassignment).
2. **Injectable roots on Request** — leaves set `CacheHome` via `t.TempDir()`
   and optional `LeafCache` / `MetricsRoot` absolute paths **outside** that home.
   Setup seeds files under those paths. Asserts only inspect `req.*` paths —
   never the developer's real user cache.
3. **No process-global env** — harness must **not** `os.Setenv` / `t.Setenv` /
   `Chdir`. Production CLI may still read `DOCTEST_CACHE_HOME`,
   `DOCTEST_LEAF_CACHE`, `DOCTEST_METRICS_ROOT`; in-process tests require an
   injectable seam.
4. **Implementer seam (preferred)** — package API such as `Scan` / `Info` /
   `Clean` (or CLI internal opts) taking `CacheHome`, extra override roots,
   `DryRun`, and writers. Wire production env → same opts; wire `Run` (or CLI
   in-process path) so `req.CacheHome` / overrides drive the op without
   process Setenv. A justified one-line update to `runWithInjectedCache` in
   this `DOCTEST.md` is acceptable if the public package API needs it; do
   **not** change leaf Assert contracts.
5. **RED phase** — `runWithInjectedCache` is plain `cli.RunWithWriters` so the
   tree compiles today; leaves fail on unknown command / missing output /
   missing deletes until the feature lands.

```go
import (
	"bytes"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/cache"
	"github.com/xhd2015/doctest/libdoc/cli"
)

// Request drives one cache CLI scenario. Leaves set Args and isolate roots on
// Request (never process env).
type Request struct {
	// Args is the full argv after the program name, e.g. ["cache"], ["cache","--clean"].
	Args []string

	// CacheHome is the injectable base cache directory (t.TempDir in leaves).
	// Product layout: CacheHome/doctest/...
	CacheHome string

	// DoctestRoot is CacheHome/doctest (set by ensureCacheHome helper).
	DoctestRoot string

	// LeafCache is an optional absolute override store root (simulates
	// DOCTEST_LEAF_CACHE when set and not under DoctestRoot). Empty = unset.
	LeafCache string

	// MetricsRoot is an optional absolute metrics base outside CacheHome
	// (simulates DOCTEST_METRICS_ROOT). Product metrics path is
	// MetricsRoot/doctest/metrics. Empty = unset.
	MetricsRoot string

	// SeededBuckets lists first-level bucket names created under DoctestRoot
	// (for asserts that check names without readdir races).
	SeededBuckets []string
}

type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Err      error
}

// Run dispatches cache/help in-process via cli.RunWithWriters (split stdout/stderr).
// CLI errors are captured into Response — never fail the harness.
// No testbin, no process Setenv/Chdir.
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	var stdout, stderr bytes.Buffer
	args := append([]string(nil), req.Args...)
	err := runWithInjectedCache(req, &stdout, &stderr, args)
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

// runWithInjectedCache is the L2 seam between the harness and product cache ops.
//
// Help / top-level / unknown-flag (no CacheHome): cli.RunWithWriters so usage and
// flag registration match the product binary.
// Info / clean with injectable CacheHome: package cache.Run with req roots so
// Parallel leaves never need process Setenv (justified seam for L2 injection).
func runWithInjectedCache(req *Request, stdout, stderr *bytes.Buffer, args []string) error {
	// Top-level help and non-cache commands always go through CLI.
	if len(args) == 0 || args[0] != "cache" {
		return cli.RunWithWriters(stdout, stderr, args)
	}
	// cache --help / -h via CLI (usage string contract).
	if len(args) > 1 && (args[1] == "--help" || args[1] == "-h") {
		return cli.RunWithWriters(stdout, stderr, args)
	}
	// Injectable roots: drive package API (info, clean, flag validation).
	if req.CacheHome != "" {
		return cache.Run(stdout, args[1:], cache.Options{
			CacheHome:   req.CacheHome,
			LeafCache:   req.LeafCache,
			MetricsRoot: req.MetricsRoot,
		})
	}
	// No injection (e.g. flags/unknown-flag): full CLI path (env-resolved roots).
	return cli.RunWithWriters(stdout, stderr, args)
}
```
