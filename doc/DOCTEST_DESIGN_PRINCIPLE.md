---
name: doctest-design-principle
description: >-
  Layer choice for tests — L1 go test, L2 doctest in-process (design center),
  L3 doctest e2e (sparse, labeled). Prefer in-process; e2e only for full
  integration. Companion to doc-spec (tree shape) and lint (parallel-safe harness).
---

--begin of skill doctest-design-principle--

# Doctest Design Principles

Practical principles for **what belongs in doctest vs go test**, and how to keep
suites fast as they grow. Companion to MECE / tree-shape review
(`doctest skill review --show`) and harness migration notes
(`doctest skill migrate --show`). This skill is about **layer choice and default
gravity**, not ASSERT frontmatter syntax or generation layout.

Show: `doctest skill design-principle --show`

---

## 1. Purpose

Doctest is optimized so that **most tests can be hierarchical scenario tests that
stay cheap**. Authors and agents should:

- Put **multi-factor public behavior** in **doctest in-process** trees (the bulk).
- Put **pure / flat edge matrices** in **go test** (a thin base).
- Put **full integration / process-boundary verification** in **doctest e2e**
  (sparse, labeled with **`e2e`** only — the public L3 identity).

Without this split, new cases default to shell-out e2e, case count grows like
unit tables, and CI wall time grows with every feature.

---

## 2. Problem: e2e becomes the default

Doctest’s strengths bias authors toward expensive full-integration paths:

| Affordance | Unintended default |
|---|---|
| Directory tree + MECE hierarchy | Scenarios feel “complete” only when end-to-end |
| Shared SETUP + shell `Run` | One expensive fixture, then every leaf reuses it |
| “Prefer doctests” in agent workflows | New leaf under an e2e root instead of in-process or `*_test.go` |
| Labels opt-in (`e2e`, `slow`, …) | Unlabeled e2e still runs in discovery; cost compounds |
| Self-hosting / nested tools | Leaves that spawn full suites multiply wall time |
| “It’s the CLI” | Help / fast-fail treated as e2e though only shallow dispatch is covered |

Result:

- **Case growth like L1, cost like L3** — many leaves, each multi-second.
- Failures surface late (after prepare + subprocess + nested work).
- False confidence that more leaves always mean more safety when many only re-hit
  the same outer path with a pure edge that belongs lower.

**Design fix:** make **doctest in-process** the default authoring surface; treat
go test and e2e as thin rails with clear jobs. **E2e is for full integration
only**, not for short paths (help, usage, early fail).

---

## 3. Three-layer model

Layers are defined by **execution model and intent**, not by “importance.”
Cost usually follows the model; labels handle exceptions.

```text
  L3 doctest e2e          ████ ~5–10% of cases
    full integration, process boundary, labeled e2e, sparse

  L2 doctest in-process   ████████████████████████ ~70–85% of cases
    ★ design center of mass — library/API and in-process CLI in same process

  L1 go test              ████ ~10–20% of cases
    pure / flat edges, package tables — no tree needed
```

| Layer | Name | Execution model | Typical SUT |
|---|---|---|---|
| **L1** | **go test** | `*_test.go` in the package | Pure helpers, parsers, local error matrices |
| **L2** | **doctest in-process** | Doctest tree; harness `Run` calls code **in the same process** | High-level APIs **and** short CLI paths via in-process CLI |
| **L3** | **doctest e2e** | **Separate process**: product binary (`testbin`/`exec`), nested `go test` / `doctest test`, real multi-step layout | Full integration contracts only |

### L2 submodes (both are in-process, not e2e)

| Submode | What `Run` does | Use for |
|---|---|---|
| **Library / API** | Call package APIs directly (e.g. `libdoc/*`, `runner.*`, `core.*`) | Policy, analyze, filter, format, multi-factor public APIs |
| **In-process CLI** | Call `cli.Run` / `cli.RunWithWriter` (same process; capturable stdout) | Help, usage, unknown flag/subcommand, skill show/list, other **short** dispatch paths |

**In-process** means the generated suite invokes code in the same address space
(including in-process CLI). **E2e** means a **process boundary** (or nested product
suite) is load-bearing in what you are proving.

Case share ≠ wall-time share. A small L3 fraction may still dominate CI minutes;
that is acceptable if L3 stays sparse and labeled, and the **bulk of cases** stay
in-process.

---

## 4. Guiding case-share ratios

Guidance for a mature package (not hard CI gates):

| Layer | Share of **cases** (approx) | Role |
|---|---|---|
| L1 go test | **10–20%** | Thin base for pure / flat edges |
| L2 doctest in-process | **70–85%** | **Main mass** — scenario hierarchy + short CLI paths |
| L3 doctest e2e | **5–10%** | Sparse **full-integration** smokes and high-value regressions |

**The whole design point:** doctest exists so most tests are **in-process
scenario trees**, not so most tests are either micro unit tables or heavy e2e.

Wall-time targets (order of magnitude, not gates):

- L1: usually tiny fraction of suite time.
- L2: most of the *runnable discovery* time, but leaves stay sub-second to low
  seconds when healthy.
- L3: large share of full/`--label-all` / `--label e2e` time; discovery should
  skip labeled leaves (including e2e) by default.

---

## 5. Author decision checklist

Before adding a case, answer **in order**:

1. **Pure / local function with many flat edges, no meaningful hierarchy?**  
   → **L1 go test** table. Stop.

2. **Short path only** (help, usage, unknown flag/subcommand, parse/validation
   fail before real work, single message + exit, shallow dispatch)?  
   → **L2 in-process** (library **or** in-process CLI such as `cli.RunWithWriter`).  
   **Not L3.** Prefer exporting a capture hook over `testbin` for these paths.

3. **Multi-factor public (or high-level package) API; shared setup; scenarios
   better as a tree?**  
   → **L2 doctest in-process** (library/API). Prefer exported entrypoints; **do not**
   shell the product binary unless the process boundary *is* the contract.

4. **Full integration** — correctness depends on a **process boundary** plus
   multi-step product behavior (real binary packaging, nested suite, multi-tool
   orchestration, real git/cwd/layout/caches that cannot be faked cheaply)?  
   → **L3 doctest e2e**, sparse, with labels (see §7).

5. **Would this leaf only re-run L3 with a flag/edge already covered at L1/L2?**  
   → **Do not add.** Keep coverage at the cheaper layer; keep one L3 smoke only if
   integration risk remains.

6. **Cost (cold, rough):**  
   - Pure micro → L1  
   - In-process ms–few seconds → L2  
   - Multi-second / nested suite → L3 + labels  

**Cheapest layer that can fail if this is wrong** wins.

---

## 6. Short paths → in-process; e2e = full integration only

### Short path (must be L2)

A case is a **short path** when it only exercises a **shallow, fast** surface:

| Short-path examples | Prefer |
|---|---|
| `--help` / usage text / subcommand lists | In-process CLI (`RunWithWriter`) |
| Unknown flag or subcommand | In-process CLI or parse API |
| Missing required arg / early validation fail | In-process CLI or library |
| Skill `--list` / `--show` content | In-process CLI |
| Pure policy on fixtures (filter, rank, format warn) | Library/API |
| Single “message + exit code” before heavy work | In-process |

**Rules:**

- Short paths **must not** be product-binary e2e leaves by default.
- Prefer **in-process CLI** (same dispatch as production, no `testbin` rebuild) or
  direct library calls.
- Unlabeled when fast; do **not** put `label: e2e` on short-path leaves.

### Full integration (L3 only)

A case is **full integration** when the **process boundary or multi-step product
path is load-bearing**:

| Full-integration examples | Stay L3 |
|---|---|
| Nested `doctest test` / full suite inside a leaf | Yes |
| Multi-tree workspace / GOCACHE / leaf-cache product path | Yes |
| Agent implement cycle with fake runner + real binary | Yes |
| Cross-module discovery, real install layout | Yes |
| Git worktree + CLI that cannot be expressed as pure filter | Maybe (prefer L2 filter + 1 L3 smoke) |

**Rules:**

- L3 is **sparse**: few leaves per capability; not N×M short-path matrices.
- L3 asserts **composition**, not combinatorics already covered at L2.
- Every L3 leaf **must** carry `label: e2e` (see §7).

### In-process CLI vs binary e2e

| | In-process CLI (L2) | Binary / nested e2e (L3) |
|---|---|---|
| Entry | `cli.Run` / `RunWithWriter` | `testbin` + `exec`, or nested product |
| Process | Same as suite | Separate process / nested suite |
| Good for | Help, fast-fail, skill show | Packaging, isolation, deep product paths |
| Labels | Usually none | **`e2e`** (required public L3 identity) |

*Is help still a “real” CLI test if in-process?* Yes for dispatch and text. Add a
rare L3 smoke only if install / `argv[0]` / binary packaging is the risk under test.

---

## 7. What each layer is optimized for

### L1 — go test (10–20%)

**Use when:**

- Pure or nearly pure functions.
- Combinatorial tables (N×M edges with no user-visible decision tree).
- Package-local helpers that do not deserve a DOCTEST root.
- Fast feedback while refactoring internals.

**Shape:** table-driven `*_test.go` next to the code. Prefer testing through a
small exported surface when possible; white-box unexported tests only when the
export surface is deliberately thin.

**Doctest role:** none for these edges. Do not explode pure tables into dozens of
directories.

### L2 — doctest in-process (70–85%) — design center

**Use when:**

- High-level **exported** APIs with **orthogonal decision factors**.
- Shared fixture/setup and outcomes that read clearly as a MECE tree.
- **Short CLI paths** via in-process CLI (help, usage, fast-fail, skill show).
- Behavior that benefits from ASSERT narrative **without** a subprocess.

**Shape:**

- One DOCTEST root per **capability**, not per micro-function.
- `Run` calls library APIs or `cli.RunWithWriter` **in-process**.
- Hierarchy ordered by **significance** (primary user-visible factor first).

**What L2 must not become:** a disguised e2e that spawns the product binary
“because SETUP already does.”

**Exhaustiveness:**

- Scenario / policy / multi-factor **public** edges → **L2 leaves** (MECE).
- Pure parse/math/string combinatorics → **L1 tables**, not fifty L2 dirs.
- Help / fast-fail matrices → **L2 in-process CLI**, not L3.

### L3 — doctest e2e (5–10%)

**Use when:**

- **Full integration** is the point: pieces work together across a process
  boundary the way a user or agent would.
- Nested test runs, multi-package prepare, agent sessions, real caches, install
  path, packaging — **not** mere help text or early fail.

**Shape:**

- **Few** leaves per capability: critical paths + known regressions.
- **Always** label with **`e2e`** (required public L3 identity). Optional program-internal
  labels (e.g. `slow`, `flaky`) may coexist; group under `e2e/` when helpful.
- Discovery skips labeled leaves by default; full CI may use `--label-all`,
  `--label e2e`, or a dedicated job.
- Prefer **one smoke per major workflow** plus targeted regressions.

**Rule of thumb:** if the assertion is fully proven by calling one exported
function, or by in-process CLI help/fast-fail, it is **not** L3.

---

## 8. Cost, labels, and discovery

### Required labels for L3

Every **true e2e** leaf (subprocess binary, nested product suite, full
integration) **must** include the run-profile label **`e2e`** in ASSERT YAML
frontmatter.

**Do not** put `label: e2e` on fixture leaves under `testdata/` or on ephemeral
fixture trees written by outer tests: nested `doctest test` discovery would skip
them and outer e2e leaves fail with “no runnable test cases.”

```yaml
---
label: e2e
explanation: nested doctest test exercises leaf-cache product path
---
```

| Label | Role |
|---|---|
| **`e2e`** | **Public L3 identity** — full integration / process-boundary. **Required on every L3 leaf.** The only publicly recognized layer label. |
| **`slow` / `flaky` / `manual`** | **Program-internal** run discipline — optional; not a substitute for `e2e` |

Multiple labels: comma-separated scalar, e.g. `label: e2e, slow`.

**Do not** put `label: e2e` on in-process leaves (including in-process CLI).  
**`heavy` is retired** as a public run-profile label; do not add it. Use **`e2e`** for
true full-integration leaves (add later when necessary). Other labels remain free-form
for program-internal filtering only.

### Discovery vs full profile

| Concern | Practice |
|---|---|
| Default discovery | Unlabeled L2 mass runs; any labeled leaf (including `e2e`) is skipped |
| Run all e2e | `doctest test --label e2e` or `--label-all` |
| How to Run | Root `DOCTEST.md` documents discovery skip and label expressions |

**`e2e` is the public layer signal.** Optional program-internal labels (e.g. `slow`)
do not define L3. An in-process leaf that is slow should stay unlabeled for layer
purposes (or use internal labels) — it is still not e2e unless a process boundary
is under test.

Nested selftest (`doctest test` / full `go test` inside a leaf) is **L3-shaped**.
Label **`e2e`**. Never use nested suite as the way to exhaust flag matrices.

---

## 9. Anti-patterns

1. **E2E by default** — new feature → new leaf under a shell-out root.  
2. **Binary e2e for help / fast-fail** — `testbin` + exec only to assert usage text
   or early exit (must be L2 in-process CLI or library).  
3. **Combinatorial e2e** — N×M flag leaves that only differ in pure parsing.  
4. **Nested suite as unit test** — leaf runs a full suite to assert one log line.  
5. **Hierarchy for pure tables** — flat twenty-edge function as twenty directories.  
6. **Unlabeled e2e** — true integration leaf without `label: e2e` (still runs in
   discovery; layer invisible to filters).  
7. **True L3 without `e2e`** — full-integration leaf missing public layer label.  
8. **Internal-only via e2e** — every internal tweak forces full stack.  
9. **Duplication without intent** — same edge at L1 and L3 with no integration-risk note.  
10. **Go-test-first for multi-factor public APIs** — starves the in-process mass that
    doctest is for (inverse of over-e2e; still wrong for the design point).

---

## 10. Positive patterns

1. **L1 owns pure edges; L2 owns policy/scenario + short CLI; L3 owns full integration.**  
2. **Default new tests to L2 in-process** for public multi-factor behavior.  
3. **Short path → `RunWithWriter` or library**; promote to binary only for deep integration.  
4. **One smoke L3 per workflow** after L1/L2 are green; always `label: e2e`.  
5. **In-process first** in SETUP/`Run`; do not copy shell-out SETUP by habit.  
6. **Extract pure cores** from CLI/handlers so L1 stays small and dense.  
7. **When adding an L3 leaf, add `label: e2e` + How to Run in the same change.**  
8. **Review question:** cheapest layer that fails if this is wrong?  
9. **When touching an old e2e-heavy tree:** demote short paths and combinatorics to
   L2; keep full-integration smokes with `label: e2e` — incremental, not big-bang.

---

## 11. Guidance for agents and TDD workflows

| Impulse | Do this instead |
|---|---|
| “Prefer doctests” | Prefer **doctest in-process** — **not** prefer e2e |
| “Add a leaf” | Default **L2**; L3 only with full-integration justification + **`label: e2e`** |
| “Add a help / usage test” | **In-process CLI** (`RunWithWriter`); unlabeled |
| “Cover all flag combos” | L1 table or L2 MECE — **not** L3 |
| “Prove the CLI works” | L2 for short paths + **one** labeled `e2e` smoke if packaging/process risk |
| “Self-test the product” | Nested full suite only as rare L3 with **`label: e2e`** |

Implementer / designer checklist (short):

1. Identify SUT: pure / short-path / multi-factor API / **full integration**.  
2. Add or extend the matching layer first.  
3. If L3: **`label: e2e`** (required public identity); optional internal labels; grouping + How to Run.  
4. Do not copy shell-out SETUP “because nearby tests do” without checking L2.  
5. Never mark in-process leaves with `e2e`.

---

## 12. Relation to other skills

| Artifact | Concern | This skill |
|---|---|---|
| MECE / significance / DSN (`doc-spec`, `review`) | **How** to structure L2/L3 trees | **Whether** the case should be L1, L2, or L3 |
| Label / run-profile audit (`review`) | Skip and layer signaling | L3 must have **`e2e`** (public); optional internal labels only |
| Parallel-safe harness (`lint`, `code-spec`) | How not to break concurrent suites | Orthogonal; layers still apply |
| Unified gen / inject (`migrate`, `code-spec`) | Harness mechanics | Orthogonal mechanics |
| Agent “prefer doctests” (`tdd`) | Often misread as prefer e2e | Prefer **in-process** doctests as the mass |

MECE still applies fully to L2 and L3 trees. It does **not** require turning pure
unit tables into directory trees.

Review skills should treat:

- missing `label: e2e` on a binary/nested leaf → **major**  
- `label: e2e` on a pure in-process leaf → **major** / mislabel  
- help/fast-fail still on `testbin` → **suggestion** or **major** (short-path rule)

---

## 13. Migrating existing heavy suites (optional, incremental)

When editing a tree that is already e2e-heavy:

1. **Classify each leaf:** short-path → L2; full integration → L3.  
2. Demote short paths to library or **in-process CLI** (e.g. `RunWithWriter`).  
3. Keep **few** full-integration smokes; set **`label: e2e`**.  
4. Drop retired `heavy` labels; add **`e2e` only when** the leaf is true full integration.  
5. Delete redundant e2e leaves that only reassert L1/L2.  
6. Aim over time toward **~70–85% in-process case share**, not overnight.

No requirement to rewrite the entire historical suite before adding new
in-process-first tests.

Follow process guards in `doctest skill lint --show` (one tree at a time,
Parallel-safe harness, no share-gaming).

---

## 14. FAQ

**Is in-process still a “real” doctest?**  
Yes. Hierarchy, SETUP, ASSERT, fixtures, and discovery all apply. Only the
execution model differs from e2e.

**Is in-process CLI the same as e2e?**  
No. Same process as the suite; use for short paths. E2e requires a separate
process or nested product suite and **`label: e2e`**.

**Should every exported function have a doctest tree?**  
No. High-level, multi-factor, scenario-worthy surfaces belong in L2. Tiny pure
exports stay L1.

**Can in-process be slow?**  
Yes (I/O, sleeps, large fixtures). Prefer fixing cost; if it must stay slow, optional
program-internal labels (e.g. `slow`) may apply. Do not call it e2e unless a process
boundary is under test.

**Does “10–20% go test” mean we under-test pure code?**  
No. Go test should still be **exhaustive** for pure edges; that slice is small in
**case count** because tables pack many edges per file, while doctest case count
is leaf-oriented. Ratios are about **where new scenario work lives**, not about
starving unit coverage.

**Why not make go test the 60–80% base (classic pyramid)?**  
That starves the product’s design point: hierarchical, reviewable scenario tests
as the **default mass**. Classic pyramids optimize for micro-units; this model
optimizes for **in-process scenario coverage** with thin pure and thin e2e rails.

**Why is `e2e` the only public L3 label?**  
**`e2e` is layer identity** (full integration / process boundary). Filters
(`--label e2e`), reviews, L3 share budgets, and demotion audits need one stable
public signal. Cost-only labels (retired `heavy`, optional internal `slow`) must
not redefine layer.

---

## 15. One-line summary

**Most tests: doctest in-process (70–85%), including short CLI paths via in-process CLI. Thin pure edges: go test (10–20%). Thin full-integration contracts: doctest e2e (5–10%), always `label: e2e` only as the public L3 identity. Prefer in-process; escalate to e2e only for real process-boundary integration.**

--end of skill doctest-design-principle--
