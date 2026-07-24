# Leaf Source-Hash Cache — Full Product (single-tree + workspace `./...` + multi-arg)

## Version
0.0.2

Specification for `github.com/xhd2015/doctest/libdoc/leafcache` and its
wiring into `doctest test` across **all product invocation shapes**.

## Layer model (L2 in-process vs L3 e2e)

Coverage-backfill layer split. **Default discovery** runs **L2 only** (unlabeled,
fast library mass). Nested product paths are **`label: heavy`** and opt-in via
`--label heavy` (or `--label-all`).

| Layer | Subtrees | Run model | Labels |
|-------|----------|-----------|--------|
| **L2 in-process** | `key/`, `store/`, `runsuite/` (nested root), `partial-package-deps/`, most of `polish/` (selective + isolation) | Call `libdoc/leafcache` APIs directly — **no** `testbin` / nested `doctest test` | unlabeled |
| **L3 e2e** | `runtime/`, `workspace/`, `cli-plan/` | Nested product binary (`testbin.Ensure` + `runtime_multi` / `runtime_once`+Bin) | **`label: e2e, heavy`** |
| **L2 help** | `polish/docs/` | `runtime_once` without Bin → `cli.RunWithWriter` | unlabeled |

| Invocation shape (L3) | Example | Branch |
|-----------------------|---------|--------|
| **Single-tree** | `doctest test ./path/to/tree` | `runtime/**` |
| **Workspace / `./...`** | `doctest test <mod>/...` | `workspace/**` |
| **Multi-arg** | `doctest test treeA treeB` | `cli-plan/**` |
| **Help docs** | `doctest test --help` | `polish/docs/**` |

Leaf-cache is **not** single-tree only. Warm skip, PutPass, disable flags
(`-count` / `-a` / `--no-leaf-cache`), and summary **Cached** apply on every
product shape above. **Cached** is the **programmatic leaf-cache skip count**
(GetPass hits), not go testcache alone — L3 leaves isolate `GOCACHE`.

| Phase | Surface | Layer | Status |
|-------|---------|-------|--------|
| **Library key/store** | `ComputeLeafKey` + `Store` | **L2** `key/**`, `store/**` | **GREEN** |
| **Partial package DAG** | multi-leaf key stability | **L2** `partial-package-deps/**` | **GREEN** |
| **Polish (keys)** | selective + isolation via ComputeLeafKey | **L2** `polish/selective/**`, `polish/isolation/**` | **GREEN** |
| **RunSuite extract** | FormatLeafIdentity / PreparePassPlan / RecordPasses | **L2** nested `runsuite/**` | **GREEN** |
| **Runtime product** | suite skip, disable flags, stream PutPass, grey dots | **L3 heavy** `runtime/**` | **GREEN** |
| **Workspace product** | multi-tree `__workspace` Cached | **L3 heavy** `workspace/**` | **GREEN** |
| **CLI multi-arg** | multi-arg `a b` Cached policy | **L3 heavy** `cli-plan/**` | **GREEN** |
| **Help docs** | `test --help` flags | **L2 in-process** `polish/docs/**` | **GREEN** |

Target share: **in-process (L2) ≥ ~60%** of leaf count; remaining e2e labeled
`heavy` so CI default discovery stays cheap.

# DSN (Domain Specific Notion)

### Participants

- **Caller / suite** — `doctest test` process that discovers leaves, may skip
  warm passes via the leaf-cache store, and reports progress.
- **Leaf spine** — ordered Go code for one leaf: root `DOCTEST.md` Go, ancestor
  `SETUP.md` Go, leaf `SETUP.md` Go, leaf `ASSERT.md` Go.
- **Local content DAG** — bottom-up hashes over spine + local import closure +
  module `go.mod`/`go.sum` + local replace modules' `go.mod` (and local sources).
- **Leaf key** — stable lowercase hex digest (algo version + GoVersion + DAG).
- **Pass store** — disk map under `$CacheHome/doctest/leaf-cache/v1` (or
  `DOCTEST_LEAF_CACHE` when set) recording **explicit** pass markers only.
- **Progress summary** — per-suite line
  `(N Run, N Pass, N Fail, N Cached)` where **Cached** counts leaves skipped
  because GetPass hit (programmatic leaf cache), not merely go testcache.

### Behaviors (library)

- **ComputeLeafKey** — local-only DAG; remote module sources not file-hashed.
- **GetPass / PutPass** — Put only on explicit pass; Get false when missing.

### Behaviors (single-tree runtime)

- **On leaf pass (stream PutPass)** — on each countable suite-leaf JSON
  `Action: pass`, suite immediately `PutPass(key)` so a mid-run interrupt
  (e.g. Ctrl-C) still leaves already-passed leaves warm for the next run.
  Fail never PutPass. End-of-run `RecordPasses` may remain as idempotent reconcile.
- **Warm skip** — before running a leaf, if leaf-cache is enabled and
  `GetPass(key)`, skip execution and count **Cached** (and treat as pass for
  totals as product defines). **Cached** is the leaf-cache product counter.
- **Grey warm progress dots (quiet + color)** — when color is on, progress `.`
  for identities in this-run leaf-cache skip set is **grey** (`\x1b[90m.\x1b[0m`).
  Executed pass stays **plain** `.`. Fail stays **red** (`\x1b[31m.\x1b[0m`).
  `-count` / `-a` / `--no-leaf-cache` disable skip → no grey leaf-cache dots.
- **Disable leaf-cache skip** when any of:
  - `-count` is set to any N (including `1`)
  - `-a` is set (also forwarded to `go test -a`)
  - `--no-leaf-cache` is set
- **Fail not stored** — failed leaves never `PutPass`; a second run still
  executes and fails (Cached stays 0 for that leaf).
- **Store I/O errors** — should not fail the suite (best-effort; optional leaf
  not required if hard to force).

### Behaviors (polish)

- **Selective invalidation** — editing one leaf's ASSERT Go changes only that
  leaf's key; sibling leaves that still GetPass stay **Cached**.
- **Local dep DAG** — editing a local package imported by a leaf invalidates
  that leaf (re-run; Cached does not claim a stale pass).
- **Tree isolation** — keys must not collide across different doctest trees that
  share the same *relative* leaf path but different absolute `TreeRoot`
  (tree identity mixed into the key or equivalent).
- **Help** — `doctest test --help` documents `-a` and `--no-leaf-cache`.

### Behaviors (RunSuite multi-prep extract)

- **FormatLeafIdentity** — tree-qualified token for skip/fail maps; distinct for
  different abs TreeRoots even when the relative leaf path is identical.
- **PreparePassPlan** — for a list of `(TreeRoot, leafRel)` refs: compute store
  keys (`map[identity]storeKey`); when skip enabled and GetPass hits, append
  identity to `Skip` (sorted). Skip empty when skip disabled.
- **RecordPasses** — PutPass for non-failed identities (or all when
  `allPassed`); never PutPass failed identities; when `!allPassed && failed==nil`,
  store nothing.
- **Store keys vs identities** — store hex keys already include abs TreeRoot;
  identities are the in-memory/env tokens for multi-tree maps (must not be bare
  relative paths alone).

### Behaviors (workspace multi-tree product — `./...` / `<mod>/...`)

- **`doctest test <mod>/...`** with multiple DOCTEST roots under one module uses
  `PrepareTree` ×N + `RunWorkspace` → `finishWorkspaceGoTest` (`__workspace`
  suite; multi-mod → `__hub` same finish path).
- **Multi-prep before go test** — call PreparePassPlan (or equivalent) for
  **all** TreePreps; pass tree-qualified skip identities into
  `DOCTEST_LEAF_CACHE_SKIP_PATHS` (generated `__wreg` / RunAll checks match
  those tokens).
- **RecordPasses after JSON** — PutPass non-failed leaves using fail identities
  from suite JSON; partial fail stores only passes.
- **Summary Cached** — leaf-cache skip count (not only go package cache);
  fresh GOCACHE in tests so Cached is programmatic leaf-cache only.
- **`-count` bypass** — any `-count=N` → SkipEnabled false → 0 Cached on
  workspace path (same as single-tree).
- **Same-relpath isolation** — warm tree-a only must not false-skip tree-b at
  the same relative path inside one workspace suite.
- **Sibling product paths** — CLI multi-arg policy lives under `cli-plan/`;
  hub multi-mod reuses the same finish path as single-module workspace (no
  separate hub leaf).

### Behaviors (CLI multi-arg product)

- **`doctest test treeA treeB`** (two or more explicit roots) respects the same
  leaf-cache enable/disable rules as single-tree and `./...` when skip is
  enabled: PutPass on pass; warm second run reports programmatic **Cached**.
- **Warm multi-arg** — two-tree all-pass fixture run twice without `-count` →
  total Cached across summary line(s) **>= number of warm leaves** (here >= 2).
- **`-count` bypass** — after a proven warm multi-arg hit, `-count=1` → 0 Cached.
- **Policy parity** — multi-arg must not ignore store / skip / summary Cached.
  Engine shape (N× `TestWithStats` vs shared SuitePlan) is free; tests assert
  **observable** Cached, not call-graph shape.
- **Already sealed elsewhere** — single-tree warm under `runtime/**`;
  `./...` workspace under `workspace/**`.

### CLI / env contract (product — all invocation shapes)

| Flag / env | Effect |
|------------|--------|
| *(default)* | leaf-cache **enabled** when `-count` unset and `-a` absent |
| `-count=N` | any N → **disable** programmatic leaf-cache skip for the invocation |
| `-a` | forward **`-a`** to `go test` (rebuild packages); **disable** leaf-cache skip |
| `--no-leaf-cache` | **disable** leaf-cache skip for this invocation |
| `DOCTEST_CACHE_HOME` | base cache home; default store at `$CacheHome/doctest/leaf-cache/v1` |
| `DOCTEST_LEAF_CACHE` | when set, **store root** for leaf pass markers (overrides default path under CacheHome) |

Runtime leaves assign a **fresh `GOCACHE` temp dir per `doctest test`
invocation** so go test result-cache never contributes `Cached` across runs.
Summary `Cached` then reflects **programmatic leaf-cache skips only**.

### Pipeline sketch

```
# Library
module + tree + leaf + goVersion -> ComputeLeafKey -> hex key
Store.PutPass / GetPass under isolated Root

# Single-tree runtime
doctest test fixture
  -> for each leaf: key = ComputeLeafKey(...)
  -> if enabled && GetPass(key): skip; Cached++; quiet+color -> grey progress .
  -> else run leaf; on Action:pass stream PutPass(key); plain progress .
  -> on Action:fail: no PutPass; red progress .
  -> end-of-run RecordPasses may reconcile (idempotent)
  -> summary (Run, Pass, Fail, Cached)

# Multi-tree workspace (./... product — not single-tree only)
doctest test <mod>/...
  -> PrepareTree each root
  -> PreparePassPlan(all trees) -> skip env (tree-qualified)
  -> go test __workspace/suite (or __hub)
  -> RecordPasses from suite JSON fails
  -> summary Cached = leaf-cache skip count

# Multi-arg CLI (same leaf-cache policy)
doctest test treeA treeB
  -> roots = union of args (may still N× TestWithStats or one SuitePlan)
  -> prepare / bind / run each root with shared store + SkipEnabled
  -> PutPass / warm Cached / -count bypass same as single + workspace
```

## Decision Tree

```
tests/leaf-cache/
├── key/                                      [L2 library ComputeLeafKey]
│   ├── stable/identical-twice
│   ├── spine/{leaf-assert-change,ancestor-setup-change}
│   ├── local-package/{imported-source-change,unrelated-source-stable}
│   ├── replace/{lib-source-change,lib-gomod-change}
│   ├── remote/source-not-hashed
│   ├── go-version/different-differs
│   └── tree-identity/different-abs-roots/    same content, different TreeRoot → keys≠
├── store/                                    [L2 library GetPass/PutPass]
│   ├── put-then-get
│   ├── missing-false
│   └── root-isolation
├── partial-package-deps/                     [L2 multi-leaf key DAG]
│   ├── edit-alone-d-two-cached/              alone/d edit → ab keys stable, d changes
│   └── edit-shared-a-one-cached/             shared/a edit → ab keys change, d stable
├── polish/                                   [mix: L2 keys + L3 help]
│   ├── selective/                            [L2]
│   │   ├── sibling-stays-cached/             leaf_a ASSERT edit → sibling key stable
│   │   └── local-dep-invalidates/            imported pkg edit → leaf key changes
│   ├── isolation/                            [L2]
│   │   └── same-relpath-two-trees/           twin TreeRoots → distinct keys
│   └── docs/                                 [L3 heavy]
│       └── test-help-mentions-flags/         test --help lists -a / no-leaf-cache
├── runsuite/                                 [L2 nested DOCTEST.md — multi-prep extract]
│   ├── identity/
│   │   ├── same-relpath-two-trees/           FormatLeafIdentity distinct across trees
│   │   └── stable-roundtrip/                 stable + distinguishes leaf rels
│   └── multi-prep/
│       ├── prepare-both-warm-skip/           both warm → 2 tree-qualified skip ids
│       ├── prepare-one-warm-only/            only A warm → skip = [idA]
│       ├── prepare-skip-disabled/            skipEnabled=false → empty Skip
│       └── record-partial-fail/              failed id not PutPass; other is
├── runtime/                                  [L3 heavy — nested doctest test]
│   ├── warm/second-run-cached/
│   ├── disable/{count,force-a,no-leaf-cache}-bypasses/
│   ├── fail-path/fail-not-stored/
│   ├── stream-pass/interrupt-partial-cached/ SIGINT after pass dots → next Cached>=1
│   └── progress-dots/
│       ├── warm-grey-dots/                   quiet+--color warm → grey progress dots
│       ├── fail-still-red/                   --color pass+fail → red fail dot
│       └── count-no-grey-cached/             warm grey then -count=1 → 0 grey dots
├── workspace/                                [L3 heavy — ./... multi-tree product]
│   ├── warm/second-run-cached/               two trees all pass; run2 Cached >= 2
│   ├── partial-fail/fail-one-others-cached/  pass+fail trees; run2 Cached >= 1, still fail
│   ├── disable/count-bypasses/               warm then -count=1 → 0 Cached
│   └── isolation/same-relpath-no-cross-skip/ warm tree-a; workspace must not false-skip tree-b
└── cli-plan/                                 [L3 heavy — multi-arg product]
    └── multi-arg/
        ├── warm/second-run-cached/           test treeA treeB twice; sum Cached >= 2
        └── disable/count-bypasses/           warm then multi-arg -count=1 → 0 Cached
```

> **Nested root:** `tests/leaf-cache/runsuite/DOCTEST.md` is a self-contained **L2**
> tree (own Request/Response/Run). **L3** product leaves (`runtime/`, `workspace/`,
> `cli-plan/`, `polish/docs/`) share the parent harness and are labeled `heavy`.

## Test Index

### P1 (sealed)

| Leaf | Scenario |
|------|----------|
| `key/stable/identical-twice` | Same inputs → identical hex keys |
| `key/spine/leaf-assert-change` | Leaf ASSERT Go change → key≠ |
| `key/spine/ancestor-setup-change` | Parent SETUP Go change → key≠ |
| `key/local-package/imported-source-change` | Imported local pkg change → key≠ |
| `key/local-package/unrelated-source-stable` | Unrelated local pkg → key stable |
| `key/replace/lib-source-change` | replace lib source → key≠ |
| `key/replace/lib-gomod-change` | replace lib go.mod → key≠ |
| `key/remote/source-not-hashed` | Remote-like file ignored |
| `key/go-version/different-differs` | Different GoVersion → keys differ |
| `store/put-then-get` | PutPass → GetPass true |
| `store/missing-false` | Missing key → false |
| `store/root-isolation` | Root A invisible under B |

### Runtime (L3 e2e single-tree product — sealed GREEN, `label: heavy`)

| Leaf | Scenario | Expect |
|------|----------|--------|
| `runtime/warm/second-run-cached` | Two default runs of a 1-pass fixture; run2 summary has Cached > 0 | **GREEN** |
| `runtime/disable/count-bypasses` | Run1 store; run2 warm Cached>0; run3 `-count=1` → `0 Cached` | **GREEN** |
| `runtime/disable/force-a-bypasses` | Run1 store; run2 warm Cached>0; run3 `-a` → `0 Cached` | **GREEN** |
| `runtime/disable/no-leaf-cache-bypasses` | Run1 store; run2 warm Cached>0; run3 `--no-leaf-cache` → `0 Cached` | **GREEN** |
| `runtime/fail-path/fail-not-stored` | Failing fixture twice; both exit ≠0 and `0 Cached` | **GREEN** |
| `runtime/stream-pass/interrupt-partial-cached` | Multi-leaf + hang; SIGINT after pass dots; unhang; run2 Cached >= 1 | **GREEN** |
| `runtime/progress-dots/warm-grey-dots` | 2-leaf warm second run with `--color` → grey progress dots >= 2 | **GREEN** |
| `runtime/progress-dots/fail-still-red` | 1-pass+1-fail `--color` → red progress dot present | **GREEN** |
| `runtime/progress-dots/count-no-grey-cached` | Warm grey then `-count=1 --color` → 0 Cached, 0 grey dots | **GREEN** |

### Polish / partial deps (mostly L2 library; docs L3)

| Leaf | Layer | Scenario | Expect |
|------|-------|----------|--------|
| `key/tree-identity/different-abs-roots` | L2 | Identical content under two abs TreeRoots → distinct keys | **GREEN** |
| `partial-package-deps/edit-alone-d-two-cached` | L2 | alone/d edit → ab keys stable, leaf-d key changes | **GREEN** |
| `partial-package-deps/edit-shared-a-one-cached` | L2 | shared/a edit → ab keys change, leaf-d stable | **GREEN** |
| `polish/selective/sibling-stays-cached` | L2 | edit leaf_a ASSERT → sibling key stable | **GREEN** |
| `polish/selective/local-dep-invalidates` | L2 | edit imported pkg → leaf key changes | **GREEN** |
| `polish/isolation/same-relpath-two-trees` | L2 | twin TreeRoots → distinct keys | **GREEN** |
| `polish/docs/test-help-mentions-flags` | **L3 heavy** | `doctest test --help` mentions `-a` and `--no-leaf-cache` | **GREEN** |

### RunSuite P1 (multi-prep extract — sealed GREEN)

| Leaf | Scenario | Expect |
|------|----------|--------|
| `runsuite/identity/same-relpath-two-trees` | FormatLeafIdentity(A,leaf) ≠ FormatLeafIdentity(B,leaf) | **GREEN** |
| `runsuite/identity/stable-roundtrip` | Format stable; different rels differ | **GREEN** |
| `runsuite/multi-prep/prepare-both-warm-skip` | both warm → 2 tree-qualified Skip ids | **GREEN** |
| `runsuite/multi-prep/prepare-one-warm-only` | only A warm → Skip = [idA] | **GREEN** |
| `runsuite/multi-prep/prepare-skip-disabled` | skipEnabled=false → empty Skip; keys still set | **GREEN** |
| `runsuite/multi-prep/record-partial-fail` | failed A not stored; B PutPass'd | **GREEN** |

### Workspace multi-tree product (L3 heavy — sealed GREEN)

| Leaf | Scenario | Expect |
|------|----------|--------|
| `workspace/warm/second-run-cached` | two-tree `/...`; run2 Cached >= 2 | **GREEN** |
| `workspace/partial-fail/fail-one-others-cached` | pass+fail trees; run2 Cached >= 1, still fail | **GREEN** |
| `workspace/disable/count-bypasses` | warm then workspace `-count=1` → 0 Cached | **GREEN** |
| `workspace/isolation/same-relpath-no-cross-skip` | warm tree-a; workspace Cached==1 and tree-b still fails | **GREEN** |

### CLI multi-arg product (L3 heavy — sealed GREEN)

| Leaf | Scenario | Expect |
|------|----------|--------|
| `cli-plan/multi-arg/warm/second-run-cached` | `test treeA treeB` twice; sum Cached >= 2 | **GREEN** |
| `cli-plan/multi-arg/disable/count-bypasses` | warm then multi-arg `-count=1` → 0 Cached | **GREEN** |

Summary **Cached** on L3 paths is the leaf-cache product counter (skip hits),
not go package cache. L2 key/store/partial/polish-keys never invoke the product binary.

## How to Run

**Discovery (default)** runs unlabeled **L2** library mass only — skips `label: heavy`.
**L3** nested product paths require `--label heavy` (or `--label-all`).

| Command | What runs |
|---------|-----------|
| `doctest test ./tests/leaf-cache/...` | L2: key, store, partial, polish keys, … |
| `doctest test --label heavy ./tests/leaf-cache/...` | L3: runtime, workspace, cli-plan, polish/docs |
| `doctest test --label-all ./tests/leaf-cache/...` | All leaves (L2 + L3) |
| Nested `runsuite/` | Always L2 (own root; no heavy labels) |

```sh
# Structure check (parent tree + nested runsuite root)
doctest vet ./tests/leaf-cache/
doctest vet ./tests/leaf-cache/runsuite/

# L2 library mass (default discovery — fast)
doctest test ./tests/leaf-cache/...
doctest test ./tests/leaf-cache/key/...
doctest test ./tests/leaf-cache/store/...
doctest test ./tests/leaf-cache/partial-package-deps/...
doctest test ./tests/leaf-cache/polish/selective/...
doctest test ./tests/leaf-cache/polish/isolation/...
doctest test ./tests/leaf-cache/runsuite/ -count=1

# L3 e2e product paths (opt-in heavy)
doctest test --label heavy ./tests/leaf-cache/...
doctest test --label heavy ./tests/leaf-cache/runtime/...
doctest test --label heavy ./tests/leaf-cache/workspace/...
doctest test --label heavy ./tests/leaf-cache/cli-plan/...
doctest test --label heavy ./tests/leaf-cache/polish/docs/...

# Full tree (library + product)
doctest test --label-all ./tests/leaf-cache/... -count=1
```

**Multi-path policy reminder:** leaf-cache applies to single-tree, workspace
`<mod>/...` / `./...`, and multi-arg `treeA treeB` with the same enable/disable
rules. Observing `N Cached` after a warm second run without `-count` proves the
product skip path for that invocation shape.

## Expected public API (library — sealed)

Package: `github.com/xhd2015/doctest/libdoc/leafcache`

- `const AlgoVersion = "v1"`
- `type KeyInput struct { ModuleRoot, TreeRoot, LeafDir, GoVersion string }`
- `func ComputeLeafKey(in KeyInput) (string, error)` — lowercase hex; abs TreeRoot mixed in
- `func NewStore(root string) (*Store, error)`
- `func (s *Store) GetPass(key string) (bool, error)` — false when missing
- `func (s *Store) PutPass(key string) error` — explicit pass only
- `func SkipEnabled(count int, force, noLeafCache bool) bool`
- `func KeyForLeaf(treeRoot, leafRel, goVersion string) (KeyInput, error)`

## Expected public API (RunSuite multi-prep extract — sealed)

Same package so single-tree `prepareLeafCache` / `recordLeafCachePasses`
and multi-tree RunSuite share one implementation:

- `func FormatLeafIdentity(treeRoot, leafRel string) string`
- `type LeafRef struct { TreeRoot, LeafRel string }`
- `type PassPlan struct { Keys map[string]string; Skip []string }`
- `func PreparePassPlan(store *Store, leaves []LeafRef, goVersion string, skipEnabled bool) (PassPlan, error)`
- `func RecordPasses(store *Store, keys map[string]string, failed map[string]bool, allPassed bool)`

In-memory models (no new on-disk format):

- `keys: map[identity]storeKey`
- `skipPaths: []identity` for GetPass hits when skip enabled
- `failed: map[identity]bool` from suite JSON

## Expected suite wiring (single-tree — sealed + stream/grey)

- On each countable suite-leaf JSON **pass**, compute key and **stream**
  `PutPass` into store under `DOCTEST_LEAF_CACHE` or
  `$CacheHome/doctest/leaf-cache/v1` (do not wait solely for end-of-run).
- Before leaf run, if skip enabled and `GetPass`, skip go test for that leaf
  and increment **Cached** in the suite summary line.
- Quiet + color: progress `.` for this-run skip-set identities is **grey**;
  executed pass plain `.`; fail **red**.
- Skip disabled when `-count` present, or `-a`, or `--no-leaf-cache`.
- Never `PutPass` on fail; partial fail stores only non-failed leaves.

## Expected suite wiring (workspace multi-tree + multi-arg — sealed product)

- In `finishWorkspaceGoTest` (and hub path): multi-prep prepare for **all**
  TreePreps before `go test`; inject skip env matching generated suite checks.
- After JSON: `RecordPasses` / `recordLeafCachePasses` with fail identities.
- Summary **Cached** is the **leaf-cache skip count** (not go testcache alone);
  product tests use a fresh `GOCACHE` so `N Cached` is programmatic leaf-cache only.
- Multi-arg `doctest test treeA treeB` respects the same SkipEnabled / PutPass /
  warm Cached policy as single-tree and workspace (see `cli-plan/**`).
- Single-tree runtime and workspace product leaves must stay GREEN together.

```go
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/cli"
	"github.com/xhd2015/doctest/libdoc/leafcache"
)

// Request selects one leaf-cache surface. Leaves set Op and related fields.
// RunSuite multi-prep lives in nested tree tests/leaf-cache/runsuite/ (own DOCTEST.md).
//
// Layer model (see ## Layer model above):
//   L2 library ops: compute_*, store_*, partial_package_keys, two_sibling_keys
//   L2 in-process:  runtime_once with Bin empty (help / short path via cli.RunWithWriter)
//   L3 e2e ops:     runtime_multi, runtime_once with Bin set (product binary)
type Request struct {
	Op string // compute_twice | compute_mutate | compute_go_versions | store_put_get | store_missing | store_isolate | compute_two_inputs | partial_package_keys | two_sibling_keys | runtime_multi | runtime_once

	// --- library fixture / store ---
	WorkDir    string
	ModuleRoot string
	TreeRoot   string
	LeafDir    string
	GoVersion  string
	GoVersionB string
	Mutation   string // leaf_assert | … | polish_edit_leaf_a | polish_edit_local_dep | polish_edit_alone_d | polish_edit_shared_a
	Flavor     string // base | replace | remote
	StoreRoot  string
	StoreRootB string
	Key        string

	// Second key input (compute_two_inputs / tree-identity).
	ModuleRootB string
	TreeRootB   string
	LeafDirB    string

	// --- runtime (doctest test subprocess) ---
	Bin        string
	Timeout    time.Duration
	FixtureDir string // mini doctest tree root (primary)
	FixtureB   string // second tree root (isolation)
	Env        []string
	Args       []string
	Args2      []string
	Args3      []string
	// MutateAfterRun: 0=none; 1=after run1 before run2; 2=after run2 before run3.
	// Uses Mutation polish_* / stream_unhang values via applyPolishMutation.
	MutateAfterRun int

	// InterruptAfterDots: when > 0, first runtime_multi invocation streams
	// stdout and sends SIGINT after this many quiet progress dots (plain or
	// colored). Proves mid-run stream PutPass durability.
	InterruptAfterDots int
}

type Response struct {
	// library key/store
	Key  string
	Key2 string
	Hit  bool
	HitB bool
	Err  string

	// runtime multi-run (up to 3 invocations)
	ExitCode  int
	Stdout    string
	Stderr    string
	ExitCode2 int
	Stdout2   string
	Stderr2   string
	ExitCode3 int
	Stdout3   string
	Stderr3   string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	resp := &Response{}

	switch req.Op {
	case "compute_twice":
		k1, err := leafcache.ComputeLeafKey(keyInput(req, req.GoVersion))
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		k2, err := leafcache.ComputeLeafKey(keyInput(req, req.GoVersion))
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		resp.Key, resp.Key2 = k1, k2
		return resp, nil

	case "compute_mutate":
		k1, err := leafcache.ComputeLeafKey(keyInput(req, req.GoVersion))
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		if err := applyMutation(t, req); err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		k2, err := leafcache.ComputeLeafKey(keyInput(req, req.GoVersion))
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		resp.Key, resp.Key2 = k1, k2
		return resp, nil

	case "compute_go_versions":
		k1, err := leafcache.ComputeLeafKey(keyInput(req, req.GoVersion))
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		k2, err := leafcache.ComputeLeafKey(keyInput(req, req.GoVersionB))
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		resp.Key, resp.Key2 = k1, k2
		return resp, nil

	case "store_put_get":
		st, err := leafcache.NewStore(req.StoreRoot)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		key := req.Key
		if key == "" {
			key, err = leafcache.ComputeLeafKey(keyInput(req, req.GoVersion))
			if err != nil {
				resp.Err = err.Error()
				return resp, err
			}
		}
		resp.Key = key
		if err := st.PutPass(key); err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		hit, err := st.GetPass(key)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		resp.Hit = hit
		return resp, nil

	case "store_missing":
		st, err := leafcache.NewStore(req.StoreRoot)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		key := req.Key
		if key == "" {
			key = "missing-key-never-written-aaaaaaaa"
		}
		resp.Key = key
		hit, err := st.GetPass(key)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		resp.Hit = hit
		return resp, nil

	case "store_isolate":
		a, err := leafcache.NewStore(req.StoreRoot)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		b, err := leafcache.NewStore(req.StoreRootB)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		key := req.Key
		if key == "" {
			key = "isolate-key-bbbbbbbbbbbbbbbbbbbbbbbbbbbb"
		}
		resp.Key = key
		if err := a.PutPass(key); err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		hitA, err := a.GetPass(key)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		hitB, err := b.GetPass(key)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		resp.Hit, resp.HitB = hitA, hitB
		return resp, nil

	case "compute_two_inputs":
		// Two independent KeyInputs (tree isolation unit test).
		k1, err := leafcache.ComputeLeafKey(leafcache.KeyInput{
			ModuleRoot: req.ModuleRoot,
			TreeRoot:   req.TreeRoot,
			LeafDir:    req.LeafDir,
			GoVersion:  req.GoVersion,
		})
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		k2, err := leafcache.ComputeLeafKey(leafcache.KeyInput{
			ModuleRoot: req.ModuleRootB,
			TreeRoot:   req.TreeRootB,
			LeafDir:    req.LeafDirB,
			GoVersion:  req.GoVersion,
		})
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		resp.Key, resp.Key2 = k1, k2
		return resp, nil

	case "partial_package_keys":
		// Multi-leaf key DAG: leaf-ab-1, leaf-ab-2 (shared a/b/c) + leaf-d (alone/d).
		// Hit = both ab leaves stable across Mutation; HitB = leaf-d stable.
		// L2 library — no product binary.
		if req.ModuleRoot == "" || req.TreeRoot == "" {
			return nil, fmt.Errorf("partial_package_keys: ModuleRoot/TreeRoot not set")
		}
		if req.Mutation == "" {
			return nil, fmt.Errorf("partial_package_keys: Mutation is empty")
		}
		ab1 := filepath.Join(req.TreeRoot, "leaf-ab-1")
		ab2 := filepath.Join(req.TreeRoot, "leaf-ab-2")
		dLeaf := filepath.Join(req.TreeRoot, "leaf-d")
		keyOf := func(leafDir string) (string, error) {
			return leafcache.ComputeLeafKey(leafcache.KeyInput{
				ModuleRoot: req.ModuleRoot,
				TreeRoot:   req.TreeRoot,
				LeafDir:    leafDir,
				GoVersion:  req.GoVersion,
			})
		}
		kAB1, err := keyOf(ab1)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		kAB2, err := keyOf(ab2)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		kD, err := keyOf(dLeaf)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		if err := applyPolishMutation(t, req); err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		kAB1b, err := keyOf(ab1)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		kAB2b, err := keyOf(ab2)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		kDb, err := keyOf(dLeaf)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		resp.Key = kAB1
		resp.Key2 = kD
		resp.Hit = kAB1 == kAB1b && kAB2 == kAB2b
		resp.HitB = kD == kDb
		return resp, nil

	case "two_sibling_keys":
		// leaf_a + leaf_b under FixtureDir; Mutation (polish_edit_leaf_a) between keys.
		// Hit = leaf_b stable; HitB = leaf_a stable. L2 library — no product binary.
		if req.FixtureDir == "" {
			return nil, fmt.Errorf("two_sibling_keys: FixtureDir not set")
		}
		if req.Mutation == "" {
			return nil, fmt.Errorf("two_sibling_keys: Mutation is empty")
		}
		// Mini trees use FixtureDir as both module and tree root (no go.mod).
		modRoot := req.ModuleRoot
		if modRoot == "" {
			modRoot = req.FixtureDir
		}
		treeRoot := req.TreeRoot
		if treeRoot == "" {
			treeRoot = req.FixtureDir
		}
		leafA := filepath.Join(req.FixtureDir, "leaf_a")
		leafB := filepath.Join(req.FixtureDir, "leaf_b")
		keyOf := func(leafDir string) (string, error) {
			return leafcache.ComputeLeafKey(leafcache.KeyInput{
				ModuleRoot: modRoot,
				TreeRoot:   treeRoot,
				LeafDir:    leafDir,
				GoVersion:  req.GoVersion,
			})
		}
		kA, err := keyOf(leafA)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		kB, err := keyOf(leafB)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		if err := applyPolishMutation(t, req); err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		kAb, err := keyOf(leafA)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		kBb, err := keyOf(leafB)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		resp.Key = kA
		resp.Key2 = kB
		resp.Hit = kB == kBb   // sibling stable
		resp.HitB = kA == kAb  // mutated leaf stable?
		return resp, nil

	case "runtime_once":
		if len(req.Args) == 0 {
			return nil, fmt.Errorf("runtime_once: req.Args is empty")
		}
		// L2 short path (help/docs): Bin empty → in-process CLI, Parallel-safe.
		if req.Bin == "" {
			var stdout bytes.Buffer
			err := cli.RunWithWriter(&stdout, req.Args)
			resp.Stdout = stdout.String()
			if err != nil {
				resp.ExitCode = 1
				resp.Stderr = err.Error()
			}
			return resp, nil
		}
		timeout := req.Timeout
		if timeout <= 0 {
			timeout = 120 * time.Second
		}
		r1, err := runDoctestCLI(t, req.Bin, req.Args, withFreshGoCache(t, req.Env), timeout)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		resp.ExitCode, resp.Stdout, resp.Stderr = r1.ExitCode, r1.Stdout, r1.Stderr
		return resp, nil

	case "runtime_multi":
		if req.Bin == "" {
			return nil, fmt.Errorf("runtime_multi: req.Bin is not set")
		}
		if len(req.Args) == 0 {
			return nil, fmt.Errorf("runtime_multi: req.Args is empty")
		}
		timeout := req.Timeout
		if timeout <= 0 {
			timeout = 180 * time.Second
		}
		// Fresh GOCACHE per invocation so go testcache cannot produce summary
		// Cached across runs; only programmatic leaf-cache skips should.
		env1 := withFreshGoCache(t, req.Env)
		var r1 *cliResult
		var err error
		if req.InterruptAfterDots > 0 {
			r1, err = runDoctestCLIInterruptAfterDots(t, req.Bin, req.Args, env1, req.InterruptAfterDots, timeout)
		} else {
			r1, err = runDoctestCLI(t, req.Bin, req.Args, env1, timeout)
		}
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		resp.ExitCode, resp.Stdout, resp.Stderr = r1.ExitCode, r1.Stdout, r1.Stderr

		if req.MutateAfterRun == 1 {
			if err := applyPolishMutation(t, req); err != nil {
				resp.Err = err.Error()
				return resp, err
			}
		}

		args2 := req.Args2
		if len(args2) == 0 {
			args2 = append([]string(nil), req.Args...)
		}
		env2 := withFreshGoCache(t, req.Env)
		r2, err := runDoctestCLI(t, req.Bin, args2, env2, timeout)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		resp.ExitCode2, resp.Stdout2, resp.Stderr2 = r2.ExitCode, r2.Stdout, r2.Stderr

		if req.MutateAfterRun == 2 {
			if err := applyPolishMutation(t, req); err != nil {
				resp.Err = err.Error()
				return resp, err
			}
		}

		if len(req.Args3) > 0 {
			env3 := withFreshGoCache(t, req.Env)
			r3, err := runDoctestCLI(t, req.Bin, req.Args3, env3, timeout)
			if err != nil {
				resp.Err = err.Error()
				return resp, err
			}
			resp.ExitCode3, resp.Stdout3, resp.Stderr3 = r3.ExitCode, r3.Stdout, r3.Stderr
		}
		return resp, nil

	default:
		return nil, fmt.Errorf("unknown Op %q", req.Op)
	}
}

func keyInput(req *Request, goVersion string) leafcache.KeyInput {
	return leafcache.KeyInput{
		ModuleRoot: req.ModuleRoot,
		TreeRoot:   req.TreeRoot,
		LeafDir:    req.LeafDir,
		GoVersion:  goVersion,
	}
}

// hexKey reports whether s looks like a non-empty lowercase hex digest.
func hexKey(s string) bool {
	if s == "" {
		return false
	}
	ok, _ := regexp.MatchString(`^[0-9a-f]+$`, s)
	return ok && len(s) >= 16
}

type cliResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

func runDoctestCLI(t *testing.T, bin string, args, extraEnv []string, timeout time.Duration) (*cliResult, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Env = append(os.Environ(), extraEnv...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	res := &cliResult{Stdout: stdout.String(), Stderr: stderr.String()}
	if err == nil {
		return res, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		res.ExitCode = exitErr.ExitCode()
		return res, nil
	}
	if ctx.Err() != nil {
		return res, ctx.Err()
	}
	return res, err
}

// runDoctestCLIInterruptAfterDots starts doctest, streams stdout, and sends
// os.Interrupt once at least afterDots quiet progress dots have been observed.
// Used by stream-pass interrupt durability leaves.
func runDoctestCLIInterruptAfterDots(t *testing.T, bin string, args, extraEnv []string, afterDots int, timeout time.Duration) (*cliResult, error) {
	t.Helper()
	if afterDots <= 0 {
		afterDots = 1
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Env = append(os.Environ(), extraEnv...)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var (
		stdoutBuf   bytes.Buffer
		stderrBuf   bytes.Buffer
		mu          sync.Mutex
		interrupted bool
	)
	stderrDone := make(chan struct{})
	go func() {
		defer close(stderrDone)
		_, _ = io.Copy(&stderrBuf, stderrPipe)
	}()

	buf := make([]byte, 512)
	for {
		n, readErr := stdoutPipe.Read(buf)
		if n > 0 {
			mu.Lock()
			stdoutBuf.Write(buf[:n])
			dots := countProgressDotsIn(stdoutBuf.String())
			if !interrupted && dots >= afterDots {
				interrupted = true
				_ = cmd.Process.Signal(os.Interrupt)
			}
			mu.Unlock()
		}
		if readErr != nil {
			break
		}
	}
	waitErr := cmd.Wait()
	<-stderrDone

	res := &cliResult{Stdout: stdoutBuf.String(), Stderr: stderrBuf.String()}
	if waitErr != nil {
		var exitErr *exec.ExitError
		if errors.As(waitErr, &exitErr) {
			res.ExitCode = exitErr.ExitCode()
		} else if ctx.Err() != nil {
			return res, ctx.Err()
		} else if !interrupted {
			return res, waitErr
		}
	}
	if !interrupted {
		return res, fmt.Errorf("never sent SIGINT: need %d progress dots; saw %d\nstdout:\n%s\nstderr:\n%s",
			afterDots, countProgressDotsIn(res.Stdout), res.Stdout, res.Stderr)
	}
	return res, nil
}

// stdoutCachedPositive reports summary Cached count > 0.
func stdoutCachedPositive(stdout string) bool {
	return cachedCount(stdout) > 0
}

// stdoutCachedZero reports an explicit 0 Cached (or no positive Cached).
func stdoutCachedZero(stdout string) bool {
	return cachedCount(stdout) == 0
}

var cachedRe = regexp.MustCompile(`(\d+)\s+Cached`)

func cachedCount(stdout string) int {
	// Prefer the last match (final suite summary).
	matches := cachedRe.FindAllStringSubmatch(stdout, -1)
	if len(matches) == 0 {
		return 0
	}
	last := matches[len(matches)-1]
	var n int
	fmt.Sscanf(last[1], "%d", &n)
	return n
}

// sumCachedCount sums every "N Cached" summary segment. Multi-arg N× trees
// emit one line per tree; workspace emits one aggregated line. Both shapes
// pass when total warm skips equal the sum.
func sumCachedCount(stdout string) int {
	matches := cachedRe.FindAllStringSubmatch(stdout, -1)
	total := 0
	for _, m := range matches {
		var n int
		fmt.Sscanf(m[1], "%d", &n)
		total += n
	}
	return total
}

// withFreshGoCache copies extraEnv and appends a unique GOCACHE temp dir.
func withFreshGoCache(t *testing.T, extraEnv []string) []string {
	t.Helper()
	out := append([]string(nil), extraEnv...)
	out = append(out, "GOCACHE="+t.TempDir())
	return out
}

// --- progress-dot helpers (quiet path; used by runtime/progress-dots/**) ---

var ansiEscapeRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

const (
	ansiGrayDot = "\x1b[90m.\x1b[0m"
	ansiRedDot  = "\x1b[31m.\x1b[0m"
)

func stripANSI(s string) string {
	return ansiEscapeRe.ReplaceAllString(s, "")
}

// progressPrefix returns stdout before the suite inline summary.
// Quiet mode often prints dots and the summary on the same line:
//
//	"..  (2 Run, 2 Pass, 0 Fail, 2 Cached) in 1s"
//
// so we cut at the summary token rather than discarding the whole line.
func progressPrefix(stdout string) string {
	// Preferred: cut at "  (N Run," (double-space before open paren).
	if loc := summaryRunRe.FindStringIndex(stdout); loc != nil {
		return stdout[:loc[0]]
	}
	// Fallback: any "(N Run,"
	if loc := summaryRunLooseRe.FindStringIndex(stdout); loc != nil {
		return stdout[:loc[0]]
	}
	// Stop before final PASS/FAIL result line if no inline summary found.
	for _, line := range strings.Split(stdout, "\n") {
		plain := strings.TrimSpace(stripANSI(line))
		if strings.HasPrefix(plain, "PASS (") || strings.HasPrefix(plain, "FAIL (") {
			idx := strings.Index(stdout, line)
			if idx >= 0 {
				return stdout[:idx]
			}
		}
	}
	return stdout
}

var (
	summaryRunRe      = regexp.MustCompile(`  \(\d+ Run,`)
	summaryRunLooseRe = regexp.MustCompile(`\(\d+ Run,`)
)

func countGrayProgressDots(stdout string) int {
	return strings.Count(progressPrefix(stdout), ansiGrayDot)
}

func countRedProgressDots(stdout string) int {
	return strings.Count(progressPrefix(stdout), ansiRedDot)
}

// countProgressDotsIn counts progress dots (plain or colored) in the progress
// region — used by interrupt-after-dots harness.
func countProgressDotsIn(stdout string) int {
	return strings.Count(stripANSI(progressPrefix(stdout)), ".")
}
```
