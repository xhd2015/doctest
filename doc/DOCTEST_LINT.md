---
name: doctest-lint
description: >-
  Hard constraints for editing doctest trees under t.Parallel: no Setenv/Chdir,
  honest L1/L2/L3 labels, one-tree demotion, no share-gaming. Companion to
  design-principle (layer choice) and code-spec (harness API).
---

--begin of skill doctest-lint--

# Doctest suite lint (agent / author checklist)

Hard constraints for editing **doctest trees and suite harnesses**. Companion to
`doctest skill design-principle --show` (what layer a case belongs on) and
`doctest skill code-spec --show` (Setup/Run/Assert API). This skill is **how not
to break the suite** while demoting or rewriting tests.

Show: `doctest skill lint --show`

Learned from real demotion work: mass rewrites, parallel-unsafe harness helpers,
label/fixture mistakes, and share-gaming via leaf deletion.

---

## 0. Process (one tree at a time)

| Rule | Detail |
|------|--------|
| **One change set = one tree (or one Run harness)** | Do not dual-mode the root, thin implementer, strip labels, and rewrite 20 Runs in one pass. |
| **Verify before next** | `doctest vet <tree>` + `doctest test <tree>/...` (+ `--label e2e` if any e2e remain). |
| **Recompute share only after green** | Inventory true e2e % after each tree; do not stack unmeasured waves. |
| **No unattended mass Run rewrite** | Scripts that replace every `func Run` are banned unless each Run body is reviewed (see §5). |
| **`git status` on the tree after edit** | Drop unexpected untracked `DOCTEST.md`, dual roots, or half-applied demotions. |

### Single-tree demotion recipe

1. **Inventory** — leaf count, nearest Run file, true e2e vs short-path.  
2. **Plan** — which leaves demote to L2; which stay L3 smokes; assert-parity map.  
3. **Edit** — that tree only.  
4. **Verify** — vet + discovery + labeled e2e for that tree.  
5. **Share** — recompute product-leaf e2e %.  
6. **Next tree** only if green.

---

## 1. Parallel-safe harness (hard fail)

Leaves run under **`t.Parallel()`** in one suite process. Process-global mutation
races other leaves.

### Forbidden in suite harness (`Run` / `Setup` for doctest trees)

| Forbidden | Why |
|-----------|-----|
| `os.Setenv` / `os.Unsetenv` / `syscall.Setenv` for leaf isolation | Process-global race |
| `os.Chdir` / `t.Chdir` for leaf workdir | Process-global race |
| Package-level mutable “session / GOCACHE / genDir” globals written per leaf | Same class of bug |
| **Package `inject*` / stash of `d.DOCTEST_*` in Setup** (e.g. `injectDoctestRoot = d.DOCTEST_ROOT`) | Same Parallel class; reintroduces free inject under new names. Prefer helpers that take `d` or path/session **strings**, or fields on `req` — full BAD/GOOD: `doctest skill code-spec --show` (**Do not re-stash d**) |
| Dual-mode that “helps” by Setenv when `Env` is non-empty | Quietly reintroduces the race |

### Required instead

| Need | Correct approach |
|------|------------------|
| Custom env or cwd for a leaf | **Subprocess** (`UseCLI` + product binary, or `exec`); set **`cmd.Env` / `cmd.Dir` only** (child process) |
| Simple CLI args, no special env/cwd | **In-process** `cli.RunWithWriter` (or library API) |
| Env/WorkDir set but subprocess path false | **Hard error** — do not fall back to Setenv/Chdir |

### In-process CLI contract

```text
if Env or WorkDir needed:
    UseCLI = true, Bin set, cmd.Env / cmd.Dir  // L3 or isolated subprocess
else:
    cli.RunWithWriter(stdout, Args)            // L2, Parallel-safe
```

Product code under test must also avoid process-global Setenv for isolation
(e.g. session id / cold GOCACHE on **opts + child `cmd.Env`**, not parent
process). Harness must obey the same rule. See
`doctest skill code-spec --show` and `doctest skill review --show`.

---

## 2. Layer honesty

Align with `doctest skill design-principle --show`.

| Layer | Execution | Labels |
|-------|-----------|--------|
| **L2 in-process** | Library API or `cli.RunWithWriter` (same process) | Usually unlabeled |
| **L3 e2e** | Product **binary** subprocess and/or load-bearing nested product suite | **`label: e2e` required** (public L3 identity) |

### Rules

- **Short path** (help, usage, unknown flag, fast-fail) → L2, never binary by default.  
- **`label: e2e` on pure L2** → mislabel (major).  
- **True L3 without `label: e2e`** → major.  
- **`heavy` is retired.** Public L3 identity is **`e2e` only**.  
- Share metrics use **execution model** (Run path), not labels alone.

---

## 3. Completeness (no share-gaming)

| Forbidden | Why |
|-----------|-----|
| Bulk-delete leaves only to hit e2e ≤ N% | Drops scenario coverage |
| Strip `e2e` labels while leaving binary `Run` | False progress |
| “In-process” that still requires parent Setenv to pass | Broken Parallel story |

### Required

- Demotion keeps the **same checkable outcomes** (Expected markers), cheaper Run.  
- Or: **L2 replacement + sparse L3 smoke** per family, with an explicit map.  
- Deleting a leaf needs: assert map, remaining smoke(s), and intentional sign-off — not an automatic side effect of a % script.

---

## 4. Fixtures

| Forbidden | Why |
|-----------|-----|
| `label: e2e` (or any skip label) on leaves under **`testdata/`** | Nested `doctest test` discovery skips them → outer leaf fails with “no runnable test cases” |
| Labeling ephemeral fixtures written only for outer tests as e2e | Same |

Fixture trees must stay **discoverable** under default (unlabeled) discovery when the
outer test expects nested leaves to run.

---

## 5. Automation guards

Before rewriting a `func Run`:

1. Read the full Run body.  
2. If it is a **multi-Op / library switch** (e.g. leaf-cache `Op=key|store|runtime`), **do not** auto-replace with a generic dual-mode CLI shell.  
3. Only auto-touch Runs that are **pure** `exec`/`testbin` wrappers with no extra product logic.  
4. After edit: `git status` on that tree; remove unexpected nested `DOCTEST.md`.  
5. Restore from git if a mass script damaged custom logic — then demote by hand.

---

## 6. Demotion recipe (checklist)

```text
[ ] Inventory this tree only (leaves, Run file, short-path vs full integration)
[ ] Plan: demote list + keep e2e smokes + assert parity
[ ] Edit harness: Parallel-safe (no Setenv/Chdir)
[ ] Edit leaves: labels match layer
[ ] doctest vet ./tests/<tree>/
[ ] doctest test ./tests/<tree>/...
[ ] doctest test --label e2e ./tests/<tree>/...   # if any e2e remain
[ ] Recompute product-leaf true-e2e %
[ ] Stop; only then start the next tree
```

---

## 7. Anti-patterns (observed)

1. **`withEnvAndDir` + `os.Setenv`/`Chdir` under Parallel** — never.  
2. **Mass dual-mode conversion of every Run** — destroyed multi-op harnesses.  
3. **Phase 0 e2e labels on fixtures under `testdata/`** — nested discovery empty.  
4. **Untracked nested `DOCTEST.md` after failed demotion** — wrong inheritance / panics.  
5. **Bulk implementer leaf deletion for share** — completeness risk; prefer planned sparse smokes + map.  
6. **All phases in parallel agents** — partial apply, hard to bisect.  
7. **Root dual-mode without forbidding Env on in-process path** — invites Setenv.  
8. **Keep-e2e Setup without `UseCLI` + `Bin`** after dual-mode default — wrong path or empty Bin.  
9. **Setup inject-stash** (`injectDoctestRoot = d.DOCTEST_ROOT`, …) — Parallel race + free-inject rename; see §1 / `code-spec`.

---

## 8. Positive patterns

1. Short path → `cli.RunWithWriter`; assert markers unchanged.  
2. Policy → package APIs in-process.  
3. True integration → sparse leaves, `label: e2e`, `cmd.Env`/`cmd.Dir` only.  
4. Product + harness: never process Setenv for session/GOCACHE isolation.  
5. Helpers take `d` or path/session strings — not package `inject*` stashed from Setup.  
6. One tree green before the next.

---

## 9. Related skills

| Skill | Concern |
|-------|---------|
| `design-principle` | L1 / L2 / L3, short-path rule, `label: e2e` |
| `migrate` | Inject, no leaf Chdir, unified gen layout |
| `code-spec` | Harness signatures, `d *session.Doctest`, Parallel API |
| `review` | Label / MECE review checklist |

---

## 10. One-line summary

**One tree at a time; never Setenv/Chdir or package inject-stash for isolation; short paths in-process; true e2e sparse + labeled; no mass Run rewrites or fixture e2e labels; verify green before the next change.**

--end of skill doctest-lint--
