# Go Best Practice Review — doctest

**Scope:** codebase structure, CLI design, flag handling, package layout  
**Method:** inspected against [go-best-practice](https://github.com/xhd2015) skill topics: `cli` (color, streaming, dry-run, skill-cli), `flags-parsing` (types, subcommand, cut, collect), `cmd-exec`, `go-embed-assets`, `kool-create`  
**Date:** 2026-08-06  
**Module:** `github.com/xhd2015/doctest`  
**Status:** review only — no implementation changes in this pass

---

## Executive summary

doctest is a mature multi-command Go CLI with a clear product surface (`test` / `build` / `vet` / `skill` / `agent` / `metrics`) and several patterns that already match go-best-practice well: less-flags on the hot paths, Shape-2 multi-skill host, plan-then-gate dry-run for cache/release, streaming list + go-test JSON progress, and solid `//go:embed` for skills and generated modules.

The main gaps are **inconsistent flag/help plumbing**, **incomplete color policy on `test`**, **raw `os/exec` instead of `xgo/support/cmd`**, and **package-layout debt** (`run/` passthrough, no `internal/`, oversized `core`/`build`, dual version sources). None of these block shipping; they raise maintenance cost and make new commands drift from the house style.

---

## What already looks good

| Area | Evidence | Topic |
|------|----------|--------|
| less-flags on primary runners | `parseTestOptions` / `parseBuildOptions` / `parseVetOptions` in `libdoc/runner/runner.go`; agent generate/implement/design in `libdoc/cli/cli.go` | `flags-parsing` |
| `**time.Duration` for go test timeout | `var timeout *time.Duration` + `Duration("-timeout,--timeout", &timeout)` so omit vs `0` stay distinct (`core.Options.Timeout`) | `flags-parsing/types` |
| Subcommand dispatch without toplevel flags | Root `run()` only checks `-h`/`--help` then switches on `args[0]` — no spurious parse of child flags | `flags-parsing/subcommand` |
| Skill CLI Shape 2 | `skill --list` / `--show` / `--install` both flag orders; `skills update`; install `--help` pass-through | `cli/skill-cli` |
| Color three-mode model | `core.ColorMode` Auto/Always/Never; mutual exclusion of `--color`/`--no-color`; TTY resolve before parallel buffers | `cli/color` |
| Streaming | `list` prints each root immediately; go test uses pipes + `-json` consumer (not full buffer-then-dump) | `cli/streaming` |
| Dry-run one pipeline | `cache.Clean(..., dryRun)`, `script/github/release` plan-then-gate, `generate.Run` with `opts.DryRun` gate | `cli/dry-run` |
| Embed for shippable content | `doc/*.md`, skill prompts, `assertmod`/`sessionmod` generated sources under `//go:embed` with tracked files | `go-embed-assets` (Layer 1) |
| Injectable CLI IO | `RunWithWriters` / `stdio` struct — no process-global stdout swap for harness parallelism | CLI UX / testability |
| Release script | uses less-flags + kool release helpers with soft-fail dry-run | `kool-create` adjacent / dry-run |

---

## Findings (by severity)

### High

#### H1. Flag parsing is a three-way mix (less-flags / manual loops / std `flag`)

**Where:**  
- less-flags: `runner` test/build/vet/list, agent generate/implement/design/with, `edit`  
- manual switch loops: `parseSkillArgs`, `cache.parseFlags`, entire `libdoc/cli/metrics.go` subcommands  
- pre-pass manual: `extractLabelFlags` for `--label` only  
- stdlib `flag`: `script/generate/embed-assert`, `embed-session`

**Why it matters:** New flags get reinvented (equals form, mutual exclusion, help). Metrics alone reimplements `--n` / `--run` / `--json` parsing three times with copy-paste equals-form handling.

**Recommended (grounded in `flags-parsing`):**

1. Standardize product CLI on less-flags only (scripts may keep `flag`).
2. Replace `extractLabelFlags` with  
   `lessflags.StringSlice("--label", &opts.LabelExprs)`  
   so `--label a --label b` and `--label=expr` work consistently.
3. Metrics: one helper or chained Parse per subcommand  
   `Bool("--json", …).Int("--n", …).String("--run", …).Help(...)`.
4. Cache:  
   `lessflags.Bool("--clean", &clean).Bool("--dry-run", &dryRun).Help("-h,--help", cacheUsage)`.

---

#### H2. `NO_COLOR` not applied for `doctest test` (list does)

**Where:**  
- `libdoc/runner/list.go` — Auto + non-empty `NO_COLOR` → Never  
- `libdoc/build/color.go` `ResolveColorMode` — TTY only; **no** `NO_COLOR`  
- `runTest` calls `ResolveColorMode` after flags → test progress ignores `NO_COLOR`

**Topic:** `cli/color` resolve order:

```text
if --color && --no-color → error
if --color               → Always
if --no-color            → Never
// Auto only:
if NO_COLOR != ""        → false
else                     → IsTerminal(stdout)
```

**Recommended:**

1. Centralize in one place (prefer `ResolveColorMode` or a shared `ResolveColor(mode, stdoutIsTTY, noColorEnv)` as in the skill).
2. Have both `list` and `test` call that single helper so Auto always honors `NO_COLOR`.
3. Document `NO_COLOR=1` in `testUsage` / `listUsage` (skill: document the env).
4. Align conflict error text with skill if you care about exact wording  
   (`--color and --no-color cannot be specified together` vs current `mutually exclusive`).

---

#### H3. Help is incomplete for core runners and several subcommands

**Topic:** `flags-parsing` / `flags-parsing/subcommand` — every level needs `-h`/`--help`; wire `lessflags.Help` next to that command’s flags.

**Gaps:**

| Command | Behavior today | Gap |
|---------|----------------|-----|
| `test` / `build` / `vet` | Help only if **first** arg is `-h`/`--help` (`runRunner`) | `doctest test ./x --help` does not hit usage; parsers have **no** `.Help(...)` |
| `edit` | Help only if first arg is help; parse has no `Help` | `doctest edit leaf --help` fails or mis-parses |
| `agent with` | First-arg help only; parse has no `Help` but handles `ErrHelp` (dead) | Mid-argv `--help` unrecognized |
| `cache` | Manual scan for help anywhere (OK-ish) but not less-flags | Divergent |
| `metrics *` | Per-subcommand `helpArgs` + `fmt.Print` | Works, but ignores `stdio` writers (see H4) |
| Skill surface | Good Shape-2 help matrix | Aligns with `cli/skill-cli` |

**Recommended:**

1. Drop first-arg-only special cases where less-flags can own help:  
   `…Help("-h,--help", usage).HelpNoExit().Parse(args)` and map `ErrHelp` → nil (as agent generate already does).
2. Keep empty-args → print usage at dispatcher level.
3. Ensure root usage continues to point at `doctest <command> --help` (already does).

---

#### H4. Metrics CLI ignores injectable stdio (`RunWithWriters` contract)

**Where:** `libdoc/cli/metrics.go` uses `fmt.Print*` / `json.NewEncoder(os.Stdout)` instead of `io.Out()` / `io.Err()`.

**Why it matters:** Nested harnesses and capture tests cannot isolate metrics output; rest of CLI carefully threads writers for parallel-safe leaves.

**Recommended:** Pass `stdio` into all `metrics*` handlers; write help and results to `io.Out()`, diagnostics to `io.Err()`.

---

### Medium

#### M1. External commands: raw `os/exec` instead of `xgo/support/cmd`

**Topic:** `cmd-exec` — prefer `github.com/xhd2015/xgo/support/cmd` for Debug printing, Dir/Env/IO chaining, and consistent capture.

**Hot spots (non-test production code):**

| Location | Use |
|----------|-----|
| `libdoc/cli/cli.go` `runAgentWith` | `exec.Command(prog, …)` inherit stdio |
| `libdoc/build/test.go` | go test JSON pipes (custom streaming — justified) |
| `libdoc/build/build.go` | `go build` + `CombinedOutput` |
| `libdoc/build/gowork.go`, `core/discover.go` | `go mod tidy` |
| `libdoc/agent/agent.go` | agent-runner exec |
| `libdoc/runner/metrics_run.go`, `metrics/resolve.go` | git remote/branch |
| `libdoc/testbin/testbin.go` | `go` build shared binary |

**Recommended:**

1. Use `cmd.Debug().Dir(…).Env(…).Run(...)` for human-facing side tools (`agent with`, git probes, tidy) so verbose/debug shows `[cmd] …`.
2. Keep custom pipe/JSON consumer for go test — document why `support/cmd` is not used (streaming events).
3. For capture-only helpers, `cmd.Output(...)` or `.Stdout(buf)` instead of ad-hoc `CombinedOutput` + error wrapping (or wrap once).
4. Note: `xgo` is already an **indirect** dependency (`go.mod`); promoting a direct require is fine if you adopt `support/cmd`.

---

#### M2. Undocumented public flag: `--go-cmd`

**Where:** Parsed in `parseTestOptions` (`String("--go-cmd", &opts.GoCmd)` + `ValidateGoCmdMode`) but **absent** from `testUsage` in `libdoc/cli/cli.go`.

**Recommended:** Document under Options: values `auto|go|xgo`, default auto (xgo when mock import detected). Or demote to `DOCTEST_DEBUG`/env if intentionally experimental.

---

#### M3. Agent `--timeout` parsed as string, not `Duration`

**Where:** `runAgentImplement` / `runAgentDesign` use `String("--timeout", &timeoutStr)` then `subagent.ParseTimeoutDuration`.

**Topic:** `flags-parsing/types` — use `Duration` (and `**time.Duration` if unset vs zero matters).

**Recommended:**  
`var timeout *time.Duration` + `lessflags.Duration("--timeout", &timeout)` if the subagent layer can accept `time.Duration`; keep custom parser only if it accepts non-Go forms beyond less-flags.

---

#### M4. `agent with` should consider Cut / clearer stop semantics

**Where:** `StopOnFirstArg()` then remaining as `prog [args…]`.

**Topic:** `flags-parsing/cut` — use Cut when the tail is an opaque foreign command and may contain flag-like tokens.

**Current shape is OK** for “own flags then prog”, but:

1. Wire `.Help("-h,--help", usage).HelpNoExit()`.
2. If users need `agent with --agent-runner=x -- --weird --flags`, document `--` / first-arg stop, or add an explicit `--exec` Cut marker for the foreign line.
3. Prefer `cmd.Debug().Env(...).Stdin(os.Stdin).Run(prog, progArgs...)` (M1).

---

#### M5. Skill list hardcodes names separate from the registry

**Where:** `runSkill` list branch hardcodes 15 names; `libdoc/spec/spec.go` `entries` is the real registry.

**Risk:** Adding a skill to `entries` without updating the list string → silent UX bug.

**Recommended:** `spec.Names()` (sorted keys of `entries`) used by both `--list` and help text generation.

---

#### M6. Package layout: public vs engine boundaries are fuzzy

**Current top-level:**

```text
cmd/doctest/     binary
run/             1-line re-export of cli.Run
libdoc/          engine (not internal/)
assert/          public assert library (also embedded via assertmod)
session/         public leaf API
doc/             embedded skill markdown
version/         VERSION embed
script/          codegen / release / hooks
tests/           doctest trees
```

**Issues:**

1. **`run/` is pure indirection** — `cmd` could import `libdoc/cli` directly, or `cli` could live under `cmd/doctest`. Extra package adds hop with no API value.
2. **`cli.Main()`** exists but process entry is `cmd` → `run.Run` → `cli.Run` — dead/alternate entry; pick one.
3. **No `internal/`** — anything under `libdoc` is importable by outsiders (`go get` module consumers). If engine packages are not a stable API, move to `internal/libdoc/…` (or document `libdoc` as public).
4. **God packages:** `libdoc/core` (~8.6k non-test lines), `libdoc/build` (~4.7k), `libdoc/runner` (~2k), `cli.go`+`metrics.go` (~1.5k). Splitting candidates: gen/assemble vs discover vs session materialize; color/summary vs test exec; metrics CLI out of `cli` into `libdoc/metrics/cli` or `cmd` subpackage.
5. **Dual public assert story is intentional but heavy:** `assert/` is the real package; `libdoc/assertmod` embeds generated copies for materialize/cache. Document this in README/dev notes so contributors do not edit `assertmod/assert.go` by hand (it is `//go:build ignore` generated).

**Topic:** `kool-create` is for greenfield scaffolds — not a fit to re-scaffold this repo. Layout fixes are evolutionary, not “kool create go-…”.

---

#### M7. Dual `VERSION.txt` + stale comment

**Where:** `cmd/doctest/VERSION.txt` and `version/VERSION.txt` (same content today).  
`version.Version()` comment says “from cmd/doctest/VERSION.txt” but embeds **local** `version/VERSION.txt`.

**Recommended:** Single source of truth (prefer `version/VERSION.txt` only); generate or copy into cmd only if install tooling needs it; fix comment.

---

#### M8. Top-level usage typo / inconsistency

**Where:** `usage` examples end with `` doc test -v ./... `` (missing “doctest”).

**Recommended:** Fix to `doctest test -v ./...` (docs-only; safe to do anytime).

---

### Low / polish

#### L1. Color policy partial elsewhere

Build path reuses test color for progress; good. Verbose go test output color is largely go’s own. No `FORCE_COLOR` / `CLICOLOR` (correct per skill).

#### L2. Dry-run coverage is good where present

- `cache --clean --dry-run` — same `PlanClean` then gate remove ✅  
- `skills update --dry-run` — delegated to skills install package ✅  
- release script — skill-shaped soft-fail dry-run ✅  
- `generate.DryRun` — same scan loop, skip write ✅  

No separate dry-run function anti-pattern found on product paths.

#### L3. Streaming: intentional buffers exist

Parallel multi-tree runs buffer per-tree then merge (`outBuf` in runner) — justified for ordered summaries. List and go-test JSON stream. Matches “buffer when you must sort/aggregate/atomic summary”.

#### L4. go:embed assets

Skill docs and prompts are committed files (not gitignored dist/) — Layer 1 compile safety is fine. No SPA hydrate story required. Codegen embed for assert/session is the right pattern for “ship generated Go into materialize”.

#### L5. Script CLIs use std `flag`

Acceptable for internal generators; optional later migration to less-flags for consistency only.

#### L6. `CollectParsedFlags` unused

No parent→child argv filter today. If you add “forward remaining flags to go test” or wrap binaries, use `flags-parsing/collect` rather than hand-filtering.

#### L7. go version pin `1.25.10`

Unusual patch pin in `go.mod`; not a go-best-practice topic, but note for contributors/CI image alignment.

---

## Topic checklist (grounded recommendations)

| Topic | Status | Action |
|-------|--------|--------|
| **flags-parsing** | Partial | Finish less-flags on metrics/cache/skill internals; StringSlice for `--label`; Help on all parsers |
| **flags-parsing/types** | Good on test timeout; weak on agent timeout | Duration for agent `--timeout` |
| **flags-parsing/subcommand** | Good root pattern | Fix mid-argv `--help` on test/build/vet/edit/with |
| **flags-parsing/cut** | N/A core; optional on `agent with` | Document StopOnFirstArg or add `--exec` Cut |
| **flags-parsing/collect** | Unused | Adopt when forwarding filtered argv |
| **cli/color** | Partial | Honor `NO_COLOR` on test via shared resolve |
| **cli/streaming** | Good | Keep list + JSON stream; keep justified buffers |
| **cli/dry-run** | Good | Keep plan-then-gate; no rewrite needed |
| **cli/skill-cli** | Good Shape 2 | DRY skill list from registry |
| **cmd-exec** | Weak | Prefer `xgo/support/cmd` for non-streaming exec |
| **go-embed-assets** | Good Layer 1 | Keep placeholders/generated embeds tracked |
| **kool-create** | N/A re-scaffold | Keep kool only for release helpers |

---

## Suggested fix order (when implementing)

1. **Color + help correctness** (H2, H3, M8) — user-visible, low risk  
2. **less-flags consolidation** (H1, M2, M3, M5) — delete manual parsers  
3. **Metrics stdio** (H4) — harness integrity  
4. **cmd wrapper** (M1, M4) — debuggability  
5. **Layout cleanup** (M6, M7) — larger refactors; do after flag/help stabilize  

---

## Package map (reference)

```text
Public API intended for leaves / installers
  assert/          output assert library
  session/         Doctest session injection API
  doc/             skill markdown (+ Content())
  version/         Version()

Binary
  cmd/doctest/     main → run.Run → libdoc/cli

Engine (currently importable)
  libdoc/cli       dispatch, skill, metrics, cache entry
  libdoc/runner    parse flags, orchestrate test/build/vet/list
  libdoc/build     generate, go test/build, color, workspace
  libdoc/core      discover, assemble, options, go-cmd, cache home
  libdoc/*         agent, designer, implementer, metrics, leafcache, …

Thin / questionable
  run/             re-export only
```

---

## Out of scope (not evaluated deeply)

- Correctness of generate/assemble algorithms  
- Performance of leaf-cache / metrics JSONL  
- Doctest tree design (see `doctest-review` skill)  
- CI workflow contents beyond layout  

---

*End of review. Implementation deferred per request.*
