---
name: doctest-migrate
description: >-
  Breaking author migration for unified suite generation: d *session.Doctest
  inject, no leaf Chdir, hierarchical gen layout, removed experiment flags.
---

--begin of skill doctest-migrate--

# Migration guide (unified suite generation)

This generation model is **breaking for authors of doctest trees** and for anyone
depending on **classic inline generation** or old **experiment flags**. Below is
a practical migration path.

Show: `doctest skill migrate --show`

Related: harness API in `doctest skill code-spec --show`; Parallel rules in
`doctest skill lint --show`.

---

## TL;DR checklist

1. **Stop using free inject vars** (`DOCTEST_ROOT` / `DOCTEST_CASE` / `DOCTEST_SESSION_ID` as package-level vars in generated code).
2. **Use `d *session.Doctest`** (second parameter after `t`) and read **`d.DOCTEST_ROOT` / `d.DOCTEST_CASE` / `d.DOCTEST_SESSION_ID`**.
3. **Do not assume leaf process cwd** is the case directory — do not rely on relative paths that only worked under leaf `Chdir`.
4. Prefer **absolute paths** built from `d.DOCTEST_CASE` / `d.DOCTEST_ROOT`.
5. Remove any scripts/CI that pass:
   - `--experiment-ref-instead-of-inline`
   - `--experiment-unified-package-per-doctest-tree`
6. Expect **unified suite layout** under gen (one suite binary per DOCTEST tree); layout probes for classic `*_test.go` per leaf will break.
7. Rebuild/reinstall the `doctest` binary, then re-run your trees cold once.
8. Refresh installed agent skills: `doctest skills update` (or reinstall).

---

## 1. Context inject: from free vars + `Chdir` → `d *session.Doctest`

### What changed

| Old world | New world |
|-----------|-----------|
| Generated tests `Chdir` into the leaf case dir | **No leaf `Chdir`**; process cwd is **undetermined** |
| Free package vars / env for root+case (legacy inject) | **`d *session.Doctest`** with fields `DOCTEST_ROOT`, `DOCTEST_CASE`, `DOCTEST_SESSION_ID` |
| Relative paths like `"./testdata"` or `"SETUP.md"` “just worked” via cwd | Use **`filepath.Join(d.DOCTEST_CASE, …)`** (or root) |

The engine populates **`d`** (authors never read free vars or getenv for inject).
Use only **`d.DOCTEST_SESSION_ID`** / `d.DOCTEST_ROOT` / `d.DOCTEST_CASE` in harness
code — never bare `DOCTEST_*` identifiers and never `os.Getenv("DOCTEST_…")`.

### Required signatures

Author Setup / Run / Assert must declare **`d *session.Doctest`** as the second
parameter after `t` (doctest validates this). Use `_ = d` if unused.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    path := filepath.Join(d.DOCTEST_CASE, "fixture.json")
    // ...
    return nil
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
    // use d.DOCTEST_ROOT / d.DOCTEST_CASE as needed
    return &Response{}, nil
}

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    // ...
}
```

Import:

```go
import "github.com/xhd2015/doctest/session"
```

(In-module trees that already depend on doctest modules get this via gen replace / session-mod as today.)

### Concrete rewrites

| Pattern to replace | Migrate to |
|--------------------|------------|
| `WorkDir: "."` / `WorkDir: "subdir"` meaning “under leaf” | `WorkDir: d.DOCTEST_CASE` or `filepath.Join(d.DOCTEST_CASE, "subdir")` |
| `os.ReadFile("ASSERT.md")` assuming leaf cwd | `os.ReadFile(filepath.Join(d.DOCTEST_CASE, "ASSERT.md"))` |
| `os.Chdir(...)` in harness for “be in case dir” | **Delete**; use absolute paths from `d` |
| Free `DOCTEST_ROOT` / `DOCTEST_CASE` in SETUP helpers | `d.DOCTEST_ROOT` / `d.DOCTEST_CASE` |
| Nested subprocesses that inherit “we chdired already” | Pass absolute paths explicitly; set `cmd.Dir` only if you truly need a child cwd |

### Helpers

- If a helper needs root/case paths, add `d *session.Doctest` (or pass the strings) — do not reintroduce package free vars.
- Package-level helpers shared across SETUP hops still work; they just should not assume cwd.
- Leaves run under **`t.Parallel()`** — never use process `Setenv`/`Chdir` for isolation (`doctest skill lint --show`).

---

## 2. Generation model: classic / experiment → hierarchical unified (default)

### What clients see under `--gen-dir` / mapping-gen

Approximate layout (module path still `testcase` for external gen trees):

```text
genRoot/
  go.mod
  <treeRel>/__droot/droot.go          # root DOCTEST/SETUP: types, Run, root setups
  <treeRel>/<parent>/setup.go         # intermediate SETUP package (if any)
  <treeRel>/<leaf>/leaf.go            # non-test leaf package (RunTestLeaf + register)
  <treeRel>/__registry/
  <treeRel>/__allleaves/
  <treeRel>/suite/suite_test.go       # ONLY package you typically `go test`
```

### Implications for your tooling

| Old assumption | New reality |
|----------------|-------------|
| One `*_test.go` per leaf, `go test ./...` over many packages | **One suite package** per tree; leaves are non-test packages |
| Grep for `TestGeneratedCase…` | Look for **`TestDoctestSuite`** and subtests (paths may encode `/` as `__`) |
| Probe `leaf_foo_test.go` | Probe **`leaf.go`** (+ suite / `__droot`) |
| Flags to enable ref/unified | **Removed** — always on (except internal-compile special case) |
| Classic “everything inlined in leaf” | **Gone as production default** |

### CLI flags removed

Delete from scripts, Makefiles, CI, docs:

```text
--experiment-ref-instead-of-inline
--experiment-unified-package-per-doctest-tree
```

There is **no** flag to restore classic production gen.

### Internal-compile exception

Trees that compile **inside** a module with **internal** package imports still use the classic assemble path (module-internal layout). Most product trees are **external gen (mapping-gen)** and get hierarchical unified. If you hit internal-compile, gen layout may still look classic for that tree only.

---

## 3. Hierarchy & authoring rules (still MECE)

- **Root** (`DOCTEST.md` / root `SETUP.md`): shared types, `Run`, root-level helpers/setup.
- **Intermediate dirs** with `SETUP.md`: become **real packages** — shared parent setup, not inlined into every leaf.
- **Leaf**: `ASSERT.md` (+ optional leaf `SETUP.md`) only; do not redefine root types / root `Run`.

Unexported symbols on root/intermediate packages are **exported by gen** (rename) so leaves can call them across packages. Prefer clear helper names; avoid locals named exactly like helpers if you ever wrote patterns like `sessionsDir := sessionsDir()` (gen only rewrites **calls**).

---

## 4. Cache / performance expectations

- Warm second run of an unchanged tree should still show **`(cached)`** for the suite package.
- Editing **leaf** SETUP/ASSERT rewrites `leaf.go` and updates the suite **leaf-source fingerprint** so result cache should **miss** (not rely on source-tree `os.Stat` — that was removed).
- Editing **intermediate/root** SETUP rewrites those packages under mapping-gen; compile graph should rebuild.
- Do **not** depend on `go clean -testcache` or forced `-count=N` as part of normal workflow.
- If you embed absolute paths in gen (case/root), changing case location implies regen — expected.

---

## 5. Suggested migration procedure (for a client repo)

1. **Upgrade** `doctest` binary / module to the release that ships unified suite gen.
2. **Search** your trees for:
   - `os.Chdir`, free `DOCTEST_ROOT` / `DOCTEST_CASE`, relative `WorkDir`, `./` assumptions
   - experiment flags above
   - layout probes for `*_test.go` / classic test names
3. **Migrate harnesses** SETUP/ASSERT/helpers to `d` + absolute paths.
4. **Run** your trees once cold, then warm:
   ```bash
   doctest test -v --label-all <your-trees>   # or your usual label policy
   doctest test -v --label-all <your-trees>
   ```
5. **Fix** failures that are path/cwd related first; then any gen-layout assertions in meta-tests.
6. **Optional:** dump gen with `--gen-dir /tmp/out` and confirm suite + `leaf.go` layout if you maintain custom tooling.
7. **Refresh skills:** `doctest skills update` so agents pick up design-principle / lint / migrate.

---

## 6. Breaking change surface (copy for release notes)

**Breaking**

- Generated leaf process no longer chdirs to the case directory.
- Inject surface is `session.Doctest` fields, not free package vars for root/case.
- Production generation is hierarchical ref packages + unified suite only.
- Experiment flags for ref/unified removed; classic default removed.
- Leaves run concurrently under `t.Parallel()` in one suite process — harness must be Parallel-safe.

**Non-goals of this migration**

- Changing assert language / v3 matching.
- Changing discovery of `DOCTEST.md` trees.
- Changing MECE / DSN authoring rules.

**Internal / advanced**

- Internal-compile trees may still use classic assemble.
- Suite leaf fingerprint may evolve; author-facing harness API should stay stable.

---

## 7. Quick before/after example

**Before (broken after upgrade):**

```go
func Setup(t *testing.T, req *Request) error {
    req.WorkDir = "." // was “leaf cwd”
    data, _ := os.ReadFile("input.txt")
    _ = data
    return nil
}
```

**After:**

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.WorkDir = d.DOCTEST_CASE
    data, err := os.ReadFile(filepath.Join(d.DOCTEST_CASE, "input.txt"))
    if err != nil {
        return err
    }
    _ = data
    return nil
}
```

---

## 8. Support / known follow-ups

- If a tree only fails after upgrade with “file not found” / wrong WorkDir → almost always **cwd migration**.
- If tooling expects many test packages → switch to **suite** + `leaf.go` layout.
- If tests race or flaky only under full suite → **Parallel-unsafe** Setenv/Chdir/globals (`doctest skill lint --show`).
- Fingerprint / intermediate FP tuning may continue without further author API breaks.

--end of skill doctest-migrate--
