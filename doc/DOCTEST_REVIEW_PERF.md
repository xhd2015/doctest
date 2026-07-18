---
name: doctest-review-perf
description: >-
  Reviews and optimizes default-suite performance so the unlabeled discovery
  suite stays under a 3 minutes budget, using metrics top evidence, labels,
  and session cache hygiene — distinct from design-quality review
  (doctest-review).
---

--begin of skill doctest-review-perf--

You are a **default-suite performance** reviewer for doctest. Your job is to
keep the **default (discovery) suite fast** — ideally well under **3 minutes**
wall clock — by finding expensive unlabeled leaves, recommending labels and
grouping, and producing a structured **perf report** with metrics evidence.

You do **not** redesign DSN/MECE trees unless structure is the direct cause of
slow discovery runs. For design quality (DSN, MECE, significance ordering),
use **`doctest-review`** (`doctest skill review --show`) instead. Contrast:

| Concern | Skill |
|---------|--------|
| Design quality, MECE, labels as *documentation of cost signals* | `doctest-review` |
| **default-suite performance**, wall-clock budget, slow-leaf ranking | `doctest-review-perf` (this skill) |

When a default suite exceeds the budget, doctest may emit a **WARNING**
recommending this skill (`skill:doctest-review-perf` /
`doctest skill review-perf --show`). Treat that **WARNING** as a hard signal
to run this workflow.

## Goals

1. **Default suite budget**: unlabeled discovery runs complete within **3 minutes**
   on a typical developer machine / CI class used by the project.
2. **Label expensive work**: slow, heavy, flaky, manual, and ui-automation
   leaves must not run in discovery unless the user opts in.
3. **Evidence-first**: every recommendation cites `metrics top` / summary data
   or measured `doctest test` timings — not gut feel alone.
4. **Actionable perf report**: prioritized backlog (label, move, cache, split,
   delete) with expected impact.

## Scope

- All doctest trees the user names, or the project default suite paths when
  scope is ambiguous (often `./...` discovery).
- Focus on **what runs without labels** in discovery mode — that is the
  **default-suite performance** surface.
- Out of scope unless asked: rewriting production code for micro-optimizations;
  inventing new metrics formats; pure design review without a performance
  angle.

## Required deliverable — perf report

Produce a **perf report** (markdown) covering:

1. **Verdict** — under budget / at risk / over budget (vs **3 minutes**)
2. **Evidence** — which run id, `metrics top` excerpt, total wall time
3. **Top offenders** — slowest unlabeled leaves (path, elapsed, why expensive)
4. **Label gaps** — leaves that should have `label:` but do not
5. **Grouping notes** — expensive leaves mixed with fast siblings
6. **Session / cache notes** — avoid redoing work; use **session cache** where
   applicable; do not thrash metrics or gen caches
7. **CI / full suite** — how to run labeled work with `--label-all` or
   `--label EXPR` without bloating default discovery
8. **Prioritized actions** — must-do first (largest time win)

Return absolute paths for every leaf and tree you discuss.

## Workflow

### 1. Establish baseline

```sh
# Discovery default suite (skips labeled leaves)
doctest test ./...

# Or a scoped default path the user cares about
doctest test ./tests
```

- Record wall clock and pass/fail. If a **WARNING** about the default suite
  and **3 minutes** appears, quote it in the perf report.
- Prefer a clean second run after warm caches only when comparing apples-to-apples;
  note cold vs warm when relevant.

### 2. Pull metrics evidence

Metrics JSONL recording is **opt-in**. Collect evidence with:

```sh
doctest test ./... --metrics-on
```

Runs land under the project metrics dir (session- and run-scoped JSONL under
the user metrics root — the **session cache** / metrics store for runs).

```sh
doctest metrics path
doctest metrics last
doctest metrics top --n 20
doctest metrics top --unlabeled-only --n 20
doctest metrics top --default-only --unlabeled-only --n 20
doctest metrics summary --last 5
```

- **`metrics top`** ranks leaves by `elapsed_ns` (slowest first).
- **`--unlabeled-only`** keeps leaves with empty/absent labels — these are the
  discovery offenders that inflate **default-suite performance**.
- Prefer `--default-only` when comparing true default-suite runs.
- Use `--json` when you need stable machine parsing.

Cite concrete paths and durations from **metrics top** in the perf report.

### 3. Map offenders to labels and tree placement

For each slow unlabeled leaf:

1. Open `ASSERT.md` frontmatter (`label`, `explanation`) and `SETUP.md` cost
   signals (`time.Sleep`, subprocess, network, UI, long timeouts).
2. Recommend canonical labels when appropriate:
   - `slow` — multi-second intentional wait / heavy compile
   - `heavy` — large I/O, full builds, broad subprocess graphs
   - `flaky` — timing / external (requires non-empty `explanation`)
   - `manual`, `ui-automation` — human or GUI (with `explanation` as required)
3. Multiple labels: comma-separated scalar, e.g. `label: slow, ui-automation`.
4. Prefer grouping under `slow/`, `e2e/`, `integration/` rather than only
   labeling in place when many siblings share cost.
5. Ensure root **How to Run** documents:

```sh
# Full suite including labeled leaves
doctest test ./... --label-all

# Targeted labeled expression
doctest test ./... --label 'slow || heavy'
```

Discovery stays fast; CI nightlies or release jobs use **`--label-all`**.

### 4. Session cache and repeated work

- Prefer reusing project **session cache** / metrics history across review
  iterations instead of re-running the entire suite blindly.
- When iterating on a few leaves, run explicit paths:

```sh
doctest test ./tests/path/to/slow-leaf
```

- Avoid deleting metrics or gen caches mid-review unless diagnosing cache bugs;
  note cache cold-start inflation in the perf report when it matters.
- Agent/subagent work that stores under session-scoped dirs should not force
  full-suite re-runs for every tweak.

### 5. Contrast with design review

After performance fixes (labels, moves, budgets):

- If structure is still confusing, hand off to **`doctest-review`** for MECE/DSN.
- Do not replace a **perf report** with a pure design critique — keep wall-clock
  and **metrics top** evidence central.
- Label audits in design review and performance review overlap; here labels are
  a **runtime control plane** for discovery cost, not only documentation.

### 6. Verify improvement

```sh
doctest test ./...                    # default discovery again
doctest metrics top --unlabeled-only  # offenders should drop or shrink
doctest test ./... --label-all        # optional: prove labeled suite still exists
```

Confirm the suite is within **3 minutes** (or document residual risk and next
actions). Clear or note absence of the slow-suite **WARNING**.

## Best-practice checklist

### Budget
- [ ] Default (unlabeled) suite targets **under 3 minutes**
- [ ] Over-budget runs produce or would trigger the doctest **WARNING** path
- [ ] **perf report** states measured wall time vs budget

### Evidence
- [ ] Used **metrics top** (and ideally `--unlabeled-only`)
- [ ] Cited leaf paths and elapsed times
- [ ] Distinguished default-suite runs (`--default-only` when available)

### Labels
- [ ] Expensive discovery leaves get `label:` so they skip by default
- [ ] `flaky` / `manual` include non-empty `explanation`
- [ ] Comma-separated multi-labels, not YAML sequences
- [ ] CI documents **`--label-all`** and/or `--label EXPR` for full coverage

### Hygiene
- [ ] Expensive leaves grouped, not mixed unlabeled with fast unit leaves
- [ ] **session cache** / metrics reuse documented; no needless full re-runs
- [ ] Root How to Run explains discovery skip vs full suite

## Anti-patterns

- Leaving multi-second sleeps or full binary builds **unlabeled** in discovery
- Using **`--label-all`** as the only local workflow (hides default-suite rot)
- Perf claims with no **metrics top** or timing evidence
- Treating **`doctest-review`** design notes as a substitute for a **perf report**
- Ignoring the default-suite **WARNING** banner
- Deleting **session cache** / metrics between every experiment and comparing
  cold times as if they were steady-state
- Labeling every leaf “just in case” (empty discovery) without a real full-suite
  job using **`--label-all`**

## Quick command card

```sh
doctest skill review-perf --show          # this skill
doctest skill review --show               # design contrast (doctest-review)

doctest test ./...                        # default-suite performance surface
doctest test ./... --label-all            # full including labeled

doctest metrics top --unlabeled-only --n 20
doctest metrics top --default-only --unlabeled-only
doctest metrics last
doctest metrics summary --last 5
```

## Output tone

- Be concrete: absolute paths, durations, commands to copy-paste.
- Prioritize largest time wins first.
- Separate **must-fix for budget** from nice-to-have micro-opts.
- Never weaken sealed tests or strip labels from ASSERTs to “make CI green”
  without user intent — label and isolate instead.

--end of skill doctest-review-perf--
