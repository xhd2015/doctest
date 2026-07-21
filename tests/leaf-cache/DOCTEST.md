# Leaf Source-Hash Cache — Key, Store (P1) + Runtime (P2) + Polish (P3)

## Version
0.0.2

Specification for `github.com/xhd2015/doctest/libdoc/leafcache` and its
wiring into `doctest test`.

| Phase | Surface | Status intent |
|-------|---------|----------------|
| **P1** | `ComputeLeafKey` + `Store` GetPass/PutPass | sealed under `key/**` (except new tree-identity), `store/**` |
| **P2** | suite skip, CLI disable flags, summary **Cached** | sealed under `runtime/**` (GREEN) |
| **P3** | selective invalidation, tree isolation, help docs | new `polish/**` + `key/tree-identity/**` |

P1/P2 ASSERTs stay sealed — P3 **adds leaves only**.

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

### Behaviors (P1 library)

- **ComputeLeafKey** — local-only DAG; remote module sources not file-hashed.
- **GetPass / PutPass** — Put only on explicit pass; Get false when missing.

### Behaviors (P2 runtime)

- **On leaf pass** — after a successful leaf, suite calls `PutPass(key)`.
- **Warm skip** — before running a leaf, if leaf-cache is enabled and
  `GetPass(key)`, skip execution and count **Cached** (and treat as pass for
  totals as product defines).
- **Disable leaf-cache skip** when any of:
  - `-count` is set to any N (including `1`)
  - `-a` is set (also forwarded to `go test -a`)
  - `--no-leaf-cache` is set
- **Fail not stored** — failed leaves never `PutPass`; a second run still
  executes and fails (Cached stays 0 for that leaf).
- **Store I/O errors** — should not fail the suite (best-effort; optional leaf
  not required if hard to force).

### Behaviors (P3 polish)

- **Selective invalidation** — editing one leaf's ASSERT Go changes only that
  leaf's key; sibling leaves that still GetPass stay **Cached**.
- **Local dep DAG** — editing a local package imported by a leaf invalidates
  that leaf (re-run; Cached does not claim a stale pass).
- **Tree isolation** — keys must not collide across different doctest trees that
  share the same *relative* leaf path but different absolute `TreeRoot`
  (tree identity mixed into the key or equivalent).
- **Help** — `doctest test --help` documents `-a` and `--no-leaf-cache`.

### CLI / env contract (P2 implementer)

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
# P1
module + tree + leaf + goVersion -> ComputeLeafKey -> hex key
Store.PutPass / GetPass under isolated Root

# P2
doctest test fixture
  -> for each leaf: key = ComputeLeafKey(...)
  -> if enabled && GetPass(key): skip; Cached++
  -> else run leaf; on pass PutPass(key)
  -> summary (Run, Pass, Fail, Cached)
```

## Decision Tree

```
tests/leaf-cache/
├── key/                                      [P1 ComputeLeafKey]
│   ├── stable/identical-twice
│   ├── spine/{leaf-assert-change,ancestor-setup-change}
│   ├── local-package/{imported-source-change,unrelated-source-stable}
│   ├── replace/{lib-source-change,lib-gomod-change}
│   ├── remote/source-not-hashed
│   └── go-version/different-differs
├── store/                                    [P1 GetPass/PutPass]
│   ├── put-then-get
│   ├── missing-false
│   └── root-isolation
├── runtime/                                  [P2 doctest test + leaf cache]
│   ├── warm/second-run-cached/
│   ├── disable/{count,force-a,no-leaf-cache}-bypasses/
│   └── fail-path/fail-not-stored/
├── key/tree-identity/                        [P3 key tree identity]
│   └── different-abs-roots/                  same content, different TreeRoot → keys≠
├── partial-package-deps/                     [partial leaf-cache across packages]
│   ├── edit-alone-d-two-cached/              edit alone/d → 2 Cached (shared leaves warm)
│   └── edit-shared-a-one-cached/             edit shared/a → 1 Cached (leaf-d warm)
└── polish/                                   [P3 runtime polish]
    ├── selective/
    │   ├── sibling-stays-cached/             edit leaf_a ASSERT → sibling still Cached
    │   └── local-dep-invalidates/            edit imported pkg → leaf re-runs (0 Cached)
    ├── isolation/
    │   └── same-relpath-two-trees/           warm treeA; treeB cold (no cross-tree hit)
    └── docs/
        └── test-help-mentions-flags/         test --help lists -a / no-leaf-cache
```

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

### P2 (new)

| Leaf | Scenario |
|------|----------|
| `runtime/warm/second-run-cached` | Two default runs of a 1-pass fixture; run2 summary has Cached > 0 |
| `runtime/disable/count-bypasses` | Run1 store; run2 warm Cached>0; run3 `-count=1` → `0 Cached` |
| `runtime/disable/force-a-bypasses` | Run1 store; run2 warm Cached>0; run3 `-a` → `0 Cached` |

| `runtime/disable/no-leaf-cache-bypasses` | Run1 store; run2 warm Cached>0; run3 `--no-leaf-cache` → `0 Cached` |
| `runtime/fail-path/fail-not-stored` | Failing fixture twice; both exit ≠0 and `0 Cached` |

### P3 (polish / isolation)

| Leaf | Scenario | Expect (observed) |
|------|----------|-------------------|
| `key/tree-identity/different-abs-roots` | Identical content under two abs TreeRoots → distinct keys | **RED** — keys currently collide (no abs TreeRoot in hash) |
| `polish/selective/sibling-stays-cached` | 2-pass tree; warm both; edit leaf_a; re-run → Cached == 1 | **GREEN** (backfill) |
| `polish/selective/local-dep-invalidates` | Pass leaf imports local pkg; warm; edit pkg; re-run → 0 Cached | **GREEN** (backfill) |
| `partial-package-deps/edit-alone-d-two-cached` | 3 leaves / 4 pkgs; run1 `-count=1`; edit alone/d; run2 → **2 Cached** | partial DAG |
| `partial-package-deps/edit-shared-a-one-cached` | same fixture; edit shared/a; run2 → **1 Cached** | multi-leaf bust |
| `polish/isolation/same-relpath-two-trees` | Warm treeA; first run treeB same relpath → 0 Cached | **RED** — treeB incorrectly Cached (same root cause) |
| `polish/docs/test-help-mentions-flags` | `doctest test --help` mentions `-a` and `--no-leaf-cache` |

## How to Run

```sh
doctest vet ./tests/leaf-cache/
doctest test ./tests/leaf-cache/ --label-all -count=1
doctest test ./tests/leaf-cache/key/...
doctest test ./tests/leaf-cache/store/...
doctest test ./tests/leaf-cache/runtime/...
doctest test ./tests/leaf-cache/polish/...
```

## Expected public API (P1 library — implementer contract)

Package: `github.com/xhd2015/doctest/libdoc/leafcache`

- `const AlgoVersion = "v1"`
- `type KeyInput struct { ModuleRoot, TreeRoot, LeafDir, GoVersion string }`
- `func ComputeLeafKey(in KeyInput) (string, error)` — lowercase hex
- `func NewStore(root string) (*Store, error)`
- `func (s *Store) GetPass(key string) (bool, error)` — false when missing
- `func (s *Store) PutPass(key string) error` — explicit pass only

## Expected suite wiring (P2 — implementer contract)

- After leaf **pass**, compute key and `PutPass` into store under
  `DOCTEST_LEAF_CACHE` or `$CacheHome/doctest/leaf-cache/v1`.
- Before leaf run, if skip enabled and `GetPass`, skip go test for that leaf
  and increment **Cached** in the suite summary line.
- Skip disabled when `-count` present, or `-a`, or `--no-leaf-cache`.
- Never `PutPass` on fail.

```go
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/leafcache"
)

// Request selects one leaf-cache surface. Leaves set Op and related fields.
type Request struct {
	Op string // compute_twice | compute_mutate | compute_go_versions | store_put_get | store_missing | store_isolate | compute_two_inputs | runtime_multi | runtime_once

	// --- P1 fixture / store ---
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

	// --- P2/P3 runtime (doctest test subprocess) ---
	Bin        string
	Timeout    time.Duration
	FixtureDir string // mini doctest tree root (primary)
	FixtureB   string // second tree root (isolation)
	Env        []string
	Args       []string
	Args2      []string
	Args3      []string
	// MutateAfterRun: 0=none; 1=after run1 before run2; 2=after run2 before run3.
	// Uses Mutation polish_* values via applyMutation / applyPolishMutation.
	MutateAfterRun int
}

type Response struct {
	// P1
	Key  string
	Key2 string
	Hit  bool
	HitB bool
	Err  string

	// P2 multi-run (up to 3 invocations)
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

func Run(t *testing.T, req *Request) (*Response, error) {
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

	case "runtime_once":
		if req.Bin == "" {
			return nil, fmt.Errorf("runtime_once: req.Bin is not set")
		}
		if len(req.Args) == 0 {
			return nil, fmt.Errorf("runtime_once: req.Args is empty")
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
		r1, err := runDoctestCLI(t, req.Bin, req.Args, env1, timeout)
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

// withFreshGoCache copies extraEnv and appends a unique GOCACHE temp dir.
func withFreshGoCache(t *testing.T, extraEnv []string) []string {
	t.Helper()
	out := append([]string(nil), extraEnv...)
	out = append(out, "GOCACHE="+t.TempDir())
	return out
}
```
