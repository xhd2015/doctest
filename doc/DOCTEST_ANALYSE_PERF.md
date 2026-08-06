---
name: doctest-analyse-perf
description: >-
  Break down doctest suite cost with built-in CLI (test --metrics-on, metrics
  last/top/phases/summary) and fill the HTML cost-breakdown template for the
  user. Cold/warm, unlabeled discovery. Not review-perf (label budget only).
---

--begin of skill doctest-analyse-perf--

# Purpose

Help the user answer **“where does the time go?”** using **doctest’s built-in
CLI** — record with `--metrics-on`, analyse with `doctest metrics …`, then
deliver a **cost breakdown** (chat and/or the **HTML report template** below).

Do **not** invent a custom profiler or invent numbers. The product surface is
the CLI; the HTML is a presentation of CLI evidence only.

| Job | Command / skill |
|-----|-----------------|
| Design quality (DSN, MECE) | `doctest skill review --show` |
| Default discovery **label budget** (~3 min) | `doctest skill review-perf --show` |
| **Measure + break down cost** (this) | `doctest skill analyse-perf --show` |

# Authority

Flag text, Env, and subcommands live in the binary. When unsure:

```sh
doctest test --help
doctest metrics --help
doctest metrics top --help
doctest metrics phases --help
```

Also: **`DOCTEST_DEBUG`** (engine-internal GODEBUG-style; listed under Env on
`doctest test --help` and section **1b** below).

# Built-in recipe (always this order)

```text
1. Record   →  doctest test <scope> --metrics-on  [+ --cold-cache / --label…]
1b. (optional) Prepare-only / host pprof via DOCTEST_DEBUG
      DOCTEST_DEBUG=bypass-go-test=1[,cpuprofile=…,memprofile=…,blockprofile=…]
      + same --metrics-on / --cold-cache as needed
2. Confirm  →  doctest metrics last
3. Pipeline →  doctest metrics phases --run last
4. Leaves   →  doctest metrics top --n 20
5. Discovery→  doctest metrics top --n 20 --unlabeled-only [--default-only]
6. Brief user (chat) and/or write HTML from the template (section 4)
7. (optional) go tool pprof on host profile files from DOCTEST_DEBUG
```

---

# 1. Record — `doctest test --metrics-on`

Metrics JSONL is **opt-in**. Without `--metrics-on`, `metrics *` has nothing
useful for this run.

```sh
# Default discovery (skips labeled leaves) — usual agent/CI surface
doctest test ./... --metrics-on

# Scoped tree or leaf
doctest test ./tests --metrics-on
doctest test ./path/to/tree --metrics-on
doctest test ./path/to/leaf --metrics-on -v

# Reproducible cold (wipe cold gen root, empty GOCACHE, force count when unset)
doctest test ./... --metrics-on --cold-cache

# Full suite including labeled leaves (different population than discovery)
doctest test ./... --metrics-on --label-all

# Expression filter
doctest test ./... --metrics-on --label e2e
```

### Cost-related flags on `test` (use with `--metrics-on`)

| Flag | Role in cost analysis |
|------|------------------------|
| **`--metrics-on`** | **Required** to record the run for later analysis |
| **`--cold-cache`** | Apples-to-apples **cold** baseline (not warm gen/testcache) |
| **`--label-all`** | Measure **full** suite (includes labeled leaves) |
| **`--label EXPR`** | Measure a **slice** of labeled work only |
| **`-count=N`** | Go test count (cold-cache forces 1 when unset) |
| **`-v`** | Verbose leaf output; noisier/slower — zoom one leaf |
| **`-cpuprofile` / `-memprofile` / …** | Forwarded to **go test** (suite binary). Use after metrics shows leaf / go_test cost |

### Host engine debug — `DOCTEST_DEBUG` (section 1b)

Engine-internal, **not** leaf harness / `d.DOCTEST_*`. Format is GODEBUG-style
comma-separated `key=value`. **Unknown keys error** (fail closed). Keys may be
**combined**. Installed at process start (all subcommands: test, vet, metrics, …).

| Key | Role in cost analysis |
|-----|------------------------|
| **`bypass-go-test=1`** | Skip host-driven **go test** after generate + workspace write+tidy. Prepare still runs. Honest summary: `BYPASS (N planned, 0 executed, go test bypassed) in …`. Isolates discover/generate wall. |
| **`cpuprofile=PATH`** | **Host** CPU profile (doctest process). Written on process exit. |
| **`memprofile=PATH`** | **Host** heap profile at exit (`GC` + write). |
| **`blockprofile=PATH`** | **Host** block profile at exit (`SetBlockProfileRate` while running). |

**Host DEBUG profiles ≠ CLI `-cpuprofile`:**

| Mechanism | What is profiled |
|-----------|------------------|
| `DOCTEST_DEBUG=cpuprofile=…` (and mem/block) | **doctest host** — discover, generate, CLI, workspace prep |
| `doctest test -cpuprofile=…` | **go test** / generated suite binary |

```sh
# Prepare-only wall (no suite compile/run)
DOCTEST_DEBUG=bypass-go-test=1 doctest test ./... --metrics-on
# → BYPASS (…); metrics phases: go_test 0 or bypassed; discover/generate dominate

# Cold prepare-only
DOCTEST_DEBUG=bypass-go-test=1 doctest test ./... --metrics-on --cold-cache

# Host CPU (+ optional mem/block) for prepare
DOCTEST_DEBUG=bypass-go-test=1,cpuprofile=/tmp/prep.pprof,memprofile=/tmp/prep.mprof,blockprofile=/tmp/prep.block \
  doctest test ./... --metrics-on
go tool pprof -http=:8080 /tmp/prep.pprof

# Host profile of any subcommand
DOCTEST_DEBUG=cpuprofile=/tmp/cpu.pprof doctest metrics last
```

Expected stderr when active (examples):

```text
doctest: DOCTEST_DEBUG cpuprofile=/tmp/prep.pprof
doctest: DOCTEST_DEBUG bypass-go-test=1 (go test will be skipped)
…
BYPASS (230 planned, 0 executed, go test bypassed) in 2.94s
```

Paths are abs-resolved from cwd; parent dirs are created. Nested child `doctest`
processes **inherit** `DOCTEST_DEBUG` — use unique profile paths or prepare-only
bypass when profiling (same path can race).

### Wall clock from the same command

`doctest test` already prints the user’s wait time. Capture:

- Pass/fail line with duration (`PASS (N/M) in …`)
- Bypass line when applicable (`BYPASS (N planned, 0 executed, go test bypassed) in …`)
- Summary: `(N Run, N Pass, N Fail, N Cached)`
- Cold-cache announce if present (`gen=… GOCACHE=… count=…`)
- `DOCTEST_DEBUG` profile banners on stderr

**Wall from this summary = “how long I waited.”**  
**Phase sums from `metrics phases` = composition of work** (can exceed wall when trees run in parallel — including **generate** under multi-tree prepare).  
With **bypass-go-test**, wall ≈ prepare; `go_test` is 0 / skipped.

---

# 2. Analyse — `doctest metrics`

```sh
doctest metrics --help
```

| Subcommand | What it answers |
|------------|-----------------|
| **`path`** | Where JSONL lives for this project |
| **`last`** | Newest run: wall, pass/skip, slowest sample |
| **`phases`** | Pipeline split: discover / generate / go_test / … + top trees by go_test |
| **`top`** | Slowest **leaves** (path + elapsed + result) |
| **`summary`** | Trend over last N runs |
| **`show`** | Full dump of one run (debug recording gaps) |
| **`prune`** | Keep newest ~30 runs |

Env: **`DOCTEST_METRICS_ROOT`** overrides metrics cache root.

## Confirm the run you just recorded

```sh
doctest metrics path
doctest metrics last
```

Human fields from `last` typically include:

```text
run_id: …
file: …
default_suite: true|false
passed: N  total: N  skipped: N
wall: …
exit_ok: true|false
leaf_count: N
slowest:
  path/to/leaf  1.2s  pass
  …
```

If `last` does not match the command you just ran (stale id / empty leaves),
re-run with `--metrics-on` and check you are in the same project cwd.

## Pipeline cost — `metrics phases`

```sh
doctest metrics phases --run last
doctest metrics phases --n 10 --json
```

Typical human output:

```text
run_id: …
suite_wall: …
phase totals (summed tree wall; may exceed suite wall when parallel):
  discover       …
  generate       …
  go_test        …
  …
top trees by go_test:
  1. path/to/tree  …
```

**How to brief the user:**

| Observation | Meaning |
|-------------|---------|
| `go_test` ≫ discover/generate | Cost is compile/link/run of generated tests, not tree walk |
| generate large | Gen/write or layout cost; compare cold vs warm |
| Phase Σ ≫ suite_wall | Parallel trees — quote **suite_wall** as wait time |
| Empty phases | Binary/run did not emit phase events — re-record with current `--metrics-on` |

## Leaf cost — `metrics top`

```sh
doctest metrics top --n 20
doctest metrics top --n 20 --unlabeled-only
doctest metrics top --n 20 --default-only --unlabeled-only
doctest metrics top --run last --json
doctest metrics top --run <run-id>
```

| Flag | Use for cost breakdown |
|------|-------------------------|
| `--n N` | How many slowest leaves |
| `--unlabeled-only` | Leaves that run in **default discovery** (no labels) |
| `--default-only` | Prefer a default-suite-shaped run when selecting |
| `--run last\|ID` | Pin which recorded run |
| `--json` | Machine-readable for scripts |

**How to brief the user:** list top paths with times; mark which are **unlabeled**
(inflate every discovery run) vs labeled (only with `--label-all` / matching
`--label`).

## Trends and dumps

```sh
doctest metrics summary --last 5
doctest metrics summary --last 5 --default-only --json
doctest metrics show              # latest
doctest metrics show <run-id>
doctest metrics prune             # retention hygiene only
```

Use **summary** when comparing several recent recordings (e.g. before/after a
change). Use **show** only when debugging empty/partial recording.

---

# 3. Deliverable — cost breakdown for the user

Always have the same evidence spine (chat **and/or** HTML):

1. **What was measured** — commands, scope, **cold vs warm**, discovery vs
   `--label-all` / `--label`, `run_id`
2. **Wall clock** — suite wall + Run / Pass / Fail / **Cached**
3. **Pipeline** — `metrics phases` (note if phase Σ > wall)
4. **Top leaves** — `metrics top`; flag **unlabeled**
5. **Discovery offenders** — `metrics top --unlabeled-only`
6. **Next 1–3 commands** — scoped, not full monorepo thrash

### Chat (always fine)

Short markdown in the conversation is enough for a quick answer.

### HTML report (use the template)

When the user wants a **saved report**, a durable archive, or a shareable
breakdown, write an HTML file by **copying the template** (section 4) and **filling every
placeholder from CLI output**.

**Fill rules:**

| Template section | Source |
|------------------|--------|
| What was measured | Exact shell commands + `metrics last` (`run_id`, `default_suite`) |
| Headline KPIs | `doctest test` summary + `metrics last` wall / leaf_count |
| Pipeline table + top trees | `metrics phases --run last` |
| Top leaves | `metrics top --n 20` |
| Discovery offenders | `metrics top --n 20 --unlabeled-only` |
| Findings | Narrative grounded only in the tables above |
| Next commands | 1–3 scoped `doctest test … --metrics-on` lines |
| Evidence appendix | Paste raw CLI stdout for last / phases / top |

- Mark cold rows with tag class `cold`, warm with `warm`.
- Unlabeled leaf rows: `class="unlabeled"` + `<span class="tag warn">unlabeled</span>`.
- Remove yellow `.fill` sample cells once real data is in.
- Never invent times; if a section has no data, write “not recorded” and say which command to re-run.

### Narrative templates (chat or Findings section)

**Pipeline-heavy**

> Suite wall **W**. `metrics phases`: go_test dominates (phase sum **P**; may
> exceed wall if trees run in parallel). Discovery/generate small. Focus on
> generated package compile/link and leaf runtime; zoom top trees with scoped
> `doctest test … --metrics-on`.

**Discovery leaf-heavy**

> Slowest **unlabeled** leaves (default discovery): `path1` **T1**, `path2` **T2**.
> These hit every unlabeled run. Consider `label: e2e` for true full integration and reserve
> full suite for `--label-all` CI. Cross-check with `review-perf` for the 3 min
> budget.

**Warm cache**

> Second run wall **W** with high **Cached** — gen/testcache hit. Compare only
> cold↔cold or warm↔warm. Re-baseline cold with `--cold-cache` if needed.

**Nested / agent trees**

> If nest/child work shows up in phases or top, report outer leaf wall vs nested
> go_test. Large nested sums mean child `doctest test` / `go test` inside leaves.

**Prepare-only / bypass**

> Suite wall **W** with `DOCTEST_DEBUG=bypass-go-test=1`. `metrics phases`:
> discover/generate dominate; go_test absent or 0. Generate **sum** may exceed
> wall (parallel trees). Host pprof (if taken) attributes engine paths
> (e.g. WriteFormattedGo / imports.Process / WalkDir) — not leaf suite runtime.
> Compare only bypass↔bypass and cold↔cold / warm↔warm.

---

# 4. HTML report template

Canonical copy in-repo: **`doc/ANALYSE_PERF_REPORT_TEMPLATE.html`**  
(keep this skill block in sync when editing the file).

Self-contained (dark theme, no external assets). Copy entire file, then replace
HTML comments `<!-- … -->` and sample rows.

```html
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>doctest cost breakdown — <!-- TITLE_SLUG --></title>
<style>
  :root {
    --bg: #0f1419;
    --panel: #1a2332;
    --border: #2d3a4d;
    --text: #e7ecf3;
    --muted: #8b9bb4;
    --accent: #5b9fd4;
    --good: #3dba7a;
    --warn: #e0a84a;
    --bad: #e05a5a;
    --mono: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    --sans: system-ui, -apple-system, "Segoe UI", Roboto, sans-serif;
  }
  * { box-sizing: border-box; }
  body {
    margin: 0;
    padding: 2rem 1.25rem 4rem;
    font-family: var(--sans);
    background: var(--bg);
    color: var(--text);
    line-height: 1.5;
    max-width: 960px;
    margin-inline: auto;
  }
  h1 { font-size: 1.55rem; font-weight: 650; margin: 0 0 0.25rem; }
  h2 {
    font-size: 1.05rem;
    font-weight: 650;
    margin: 2rem 0 0.75rem;
    padding-bottom: 0.35rem;
    border-bottom: 1px solid var(--border);
    color: var(--accent);
  }
  h3 { font-size: 0.95rem; margin: 1.25rem 0 0.5rem; color: var(--muted); font-weight: 600; }
  p, li { color: var(--text); }
  .muted { color: var(--muted); }
  .sub { color: var(--muted); font-size: 0.9rem; margin: 0 0 1.25rem; }
  .panel {
    background: var(--panel);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 1rem 1.1rem;
    margin: 0.75rem 0;
  }
  .kpis {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
    gap: 0.65rem;
    margin: 0.75rem 0 1rem;
  }
  .kpi {
    background: var(--panel);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 0.75rem 0.9rem;
  }
  .kpi .label { font-size: 0.72rem; text-transform: uppercase; letter-spacing: 0.04em; color: var(--muted); }
  .kpi .value { font-size: 1.25rem; font-weight: 650; font-family: var(--mono); margin-top: 0.2rem; }
  .kpi .hint { font-size: 0.75rem; color: var(--muted); margin-top: 0.15rem; }
  table {
    width: 100%;
    border-collapse: collapse;
    font-size: 0.88rem;
    margin: 0.5rem 0 0.25rem;
  }
  th, td {
    text-align: left;
    padding: 0.45rem 0.55rem;
    border-bottom: 1px solid var(--border);
    vertical-align: top;
  }
  th { color: var(--muted); font-weight: 600; font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.03em; }
  td.path, code, pre { font-family: var(--mono); }
  td.num, th.num { text-align: right; font-family: var(--mono); white-space: nowrap; }
  tr.unlabeled td.path { color: var(--warn); }
  .tag {
    display: inline-block;
    font-size: 0.7rem;
    font-weight: 600;
    padding: 0.1rem 0.4rem;
    border-radius: 4px;
    border: 1px solid var(--border);
    color: var(--muted);
    vertical-align: middle;
  }
  .tag.cold { color: var(--accent); border-color: #3a5f80; }
  .tag.warm { color: var(--good); border-color: #2a6b4a; }
  .tag.warn { color: var(--warn); border-color: #8a6a2a; }
  .tag.bad { color: var(--bad); border-color: #8a3a3a; }
  pre {
    background: #0a0e14;
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 0.75rem 0.9rem;
    overflow-x: auto;
    font-size: 0.78rem;
    line-height: 1.4;
    color: #c5d0e0;
    margin: 0.4rem 0;
  }
  ul.cmds { padding-left: 1.2rem; }
  ul.cmds li { margin: 0.35rem 0; }
  footer {
    margin-top: 2.5rem;
    padding-top: 1rem;
    border-top: 1px solid var(--border);
    font-size: 0.8rem;
    color: var(--muted);
  }
  .note {
    font-size: 0.85rem;
    color: var(--muted);
    border-left: 3px solid var(--border);
    padding-left: 0.75rem;
    margin: 0.75rem 0;
  }
  .fill {
    background: rgba(224, 168, 74, 0.12);
    border-radius: 3px;
    padding: 0 0.15rem;
  }
</style>
</head>
<body>

<header>
  <h1>doctest cost breakdown</h1>
  <p class="sub">
    <!-- PROJECT_OR_SCOPE --> ·
    <span class="tag <!-- COLD_OR_WARM_CLASS -->"><!-- COLD_OR_WARM --></span>
    <span class="tag"><!-- DISCOVERY_OR_LABEL_MODE --></span>
    · <!-- DATE_ISO -->
  </p>
</header>

<section>
  <h2>1. What was measured</h2>
  <div class="panel">
    <table>
      <tr><th>Field</th><th>Value</th></tr>
      <tr><td>Scope</td><td class="path"><!-- SCOPE e.g. ./... --></td></tr>
      <tr><td>Commands</td><td><pre><!-- COMMANDS_BLOCK --></pre></td></tr>
      <tr><td>run_id</td><td class="path"><!-- RUN_ID from metrics last --></td></tr>
      <tr><td>default_suite</td><td><!-- true|false --></td></tr>
      <tr><td>doctest / go</td><td class="muted"><!-- optional versions --></td></tr>
      <tr><td>Notes</td><td><!-- cold-cache announce, GOCACHE, gen dir, count; DOCTEST_DEBUG=bypass / host profile paths --></td></tr>
    </table>
  </div>
</section>

<section>
  <h2>2. Headline (wall clock)</h2>
  <p class="note">
    <strong>Wall</strong> = how long the user waited (<code>doctest test</code> summary /
    <code>metrics last</code> / <code>phases</code> suite_wall).
    Phase sums may exceed wall when trees run in parallel.
  </p>
  <div class="kpis">
    <div class="kpi">
      <div class="label">Suite wall</div>
      <div class="value"><!-- WALL e.g. 42.1s --></div>
      <div class="hint"><!-- cold|warm --></div>
    </div>
    <div class="kpi">
      <div class="label">Run / Pass / Fail</div>
      <div class="value"><!-- R / P / F --></div>
      <div class="hint">from test summary</div>
    </div>
    <div class="kpi">
      <div class="label">Cached</div>
      <div class="value"><!-- CACHED_N --></div>
      <div class="hint">go testcache hits</div>
    </div>
    <div class="kpi">
      <div class="label">Leaves (metrics)</div>
      <div class="value"><!-- LEAF_COUNT --></div>
      <div class="hint">leaf_end events</div>
    </div>
  </div>
</section>

<section>
  <h2>3. Pipeline — <code>metrics phases</code></h2>
  <div class="panel">
    <table>
      <thead>
        <tr><th>Phase</th><th class="num">Summed tree wall</th><th>Share / note</th></tr>
      </thead>
      <tbody>
        <!-- PHASE_ROWS from metrics phases -->
        <tr><td class="path fill">PHASE</td><td class="num fill">TIME</td><td class="muted fill">NOTE</td></tr>
      </tbody>
    </table>
    <p class="muted" style="margin:0.6rem 0 0;font-size:0.85rem">
      Phase total Σ: <!-- PHASE_SUM --> · suite_wall: <!-- WALL --> ·
      <!-- PARALLEL_NOTE -->
    </p>
  </div>
  <h3>Top trees by go_test</h3>
  <div class="panel">
    <table>
      <thead>
        <tr><th>#</th><th>Tree</th><th class="num">go_test</th></tr>
      </thead>
      <tbody>
        <!-- TOP_TREE_ROWS -->
        <tr><td class="num fill">1</td><td class="path fill">TREE</td><td class="num fill">TIME</td></tr>
      </tbody>
    </table>
  </div>
</section>

<section>
  <h2>4. Top leaves — <code>metrics top</code></h2>
  <div class="panel">
    <table>
      <thead>
        <tr>
          <th>#</th><th>Path</th><th class="num">Elapsed</th><th>Result</th><th>Labels</th>
        </tr>
      </thead>
      <tbody>
        <!-- TOP_LEAF_ROWS; class="unlabeled" when no labels -->
        <tr class="unlabeled">
          <td class="num fill">1</td>
          <td class="path fill">LEAF_PATH</td>
          <td class="num fill">TIME</td>
          <td class="fill">pass</td>
          <td><span class="tag warn">unlabeled</span></td>
        </tr>
      </tbody>
    </table>
  </div>
  <h3>Discovery offenders — <code>metrics top --unlabeled-only</code></h3>
  <div class="panel">
    <table>
      <thead>
        <tr><th>#</th><th>Path</th><th class="num">Elapsed</th><th>Impact</th></tr>
      </thead>
      <tbody>
        <!-- UNLABELED_ROWS -->
        <tr>
          <td class="num fill">1</td>
          <td class="path fill">LEAF_PATH</td>
          <td class="num fill">TIME</td>
          <td class="muted fill">hits every default discovery run</td>
        </tr>
      </tbody>
    </table>
  </div>
</section>

<section>
  <h2>5. Findings</h2>
  <div class="panel">
    <ul>
      <!-- FINDINGS -->
      <li class="fill">Dominant cost is … (cite phase + wall).</li>
      <li class="fill">Slowest leaves: …</li>
      <li class="fill">Discovery exposure: …</li>
    </ul>
  </div>
</section>

<section>
  <h2>6. Recommended next commands</h2>
  <div class="panel">
    <ul class="cmds">
      <li><pre>doctest test ./path/to/slow/tree --metrics-on</pre></li>
      <li><pre>doctest test ./path/to/leaf --metrics-on -v</pre></li>
    </ul>
  </div>
</section>

<section>
  <h2>7. Evidence appendix (CLI paste)</h2>
  <h3><code>doctest metrics last</code></h3>
  <pre><!-- PASTE metrics last --></pre>
  <h3><code>doctest metrics phases --run last</code></h3>
  <pre><!-- PASTE metrics phases --></pre>
  <h3><code>doctest metrics top --n 20</code></h3>
  <pre><!-- PASTE metrics top --></pre>
  <h3><code>doctest metrics top --n 20 --unlabeled-only</code></h3>
  <pre><!-- PASTE metrics top unlabeled --></pre>
</section>

<footer>
  Generated for analyse-perf · evidence from
  <code>doctest test --metrics-on</code> + <code>doctest metrics</code>.
  Replace every placeholder from CLI output; do not invent numbers.
</footer>

</body>
</html>
```

---

# Wrong → correct

| Wrong | Correct |
|-------|---------|
| Analyse without recording | `doctest test … --metrics-on` first |
| Quote phase sum as “user wait” | Wall from `test` / suite_wall; phases for composition |
| Cold vs warm as regression | Label cold/warm; use `--cold-cache` for cold |
| Profile whole monorepo first | `metrics top` → scoped tree → optional pprof |
| Rank discovery cost ignoring labels | `--unlabeled-only` / `--default-only` |
| Free-form HTML with invented numbers | Fill **this** template from CLI; paste evidence appendix |
| HTML before measuring | Record + `metrics *` first, then fill template |
| `go run -cpuprofile=… ./cmd/doctest` | Build/install binary + `DOCTEST_DEBUG=cpuprofile=…` (or OS sample) |
| CLI `-cpuprofile` to profile generate/discover | Host `DOCTEST_DEBUG=cpuprofile=` |
| Quote generate phase **sum** as prepare wait | Suite wall; generate sum ≫ wall when trees prepare in parallel |
| Same host profile path for nested full suite | Unique paths, or profile with `bypass-go-test=1` |

# Pitfalls

- Metrics **off by default** — missing `--metrics-on` → empty/stale analysis  
- Stale run: always `metrics last` after the command you care about  
- Parallel trees: go_test **and** generate phase sums can be several× wall  
- `--label-all` measures a **different** suite than default discovery  
- `--cold-cache` changes gen/GOCACHE — say so in the breakdown  
- Nested selftest binaries: wipe selftest-bin only when measuring **this**
  doctest binary’s own tests  
- Leaving `.fill` sample rows in the published HTML  
- Nested child `doctest` **inherits** `DOCTEST_DEBUG` — same `cpuprofile` path can
  race; prefer unique paths or prepare-only bypass  
- `bypass-go-test` is **engine debug**, not a product feature / CI default  
- Host `memprofile` is heap **at exit**, not continuous alloc sampling  

# Related

- `doc/ANALYSE_PERF_REPORT_TEMPLATE.html` — same HTML template as a file  
- `doctest skill review-perf --show` — labels + default-suite budget  
- `doctest metrics --help` / `doctest test --help` — flag + **Env `DOCTEST_DEBUG`** authority  
- `doc/GO_TEST_CACHE.md` — gen/mtime/chdir notes  
- Cache digs: `script/debug/bug-repro/go-test-cache-stale-even-code-changed/`

--end of skill doctest-analyse-perf--
