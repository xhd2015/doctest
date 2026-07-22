# Scenario

**Feature**: leaf source-hash cache — L2 library mass + L3 nested product paths

```
# L2 library (unlabeled; default discovery)
module + tree + leaf + goVersion
  -> spine Go + local pkgs + go.mod/sum + local replace
  -> hex key (abs TreeRoot mixed in)
Caller -> Store.PutPass(key) -> GetPass(key)=true
partial multi-leaf / polish selective keys via ComputeLeafKey only

# L2 RunSuite extract (nested DOCTEST.md)
FormatLeafIdentity / PreparePassPlan / RecordPasses
  -> tree-qualified skip/fail map identities

# L3 e2e (label: heavy) — nested product binary
doctest test fixture / <mod>/... / tree-a tree-b
  -> PutPass / warm Cached / -count|-a|--no-leaf-cache bypass
  -> stream PutPass + grey progress dots (runtime/**)
```

## Preconditions

- Package under test: `github.com/xhd2015/doctest/libdoc/leafcache`.
- **L2** leaves (`key/`, `store/`, `partial-package-deps/`, `polish/selective|isolation/`):
  set `req.Op` to library ops; **no** `testbin`.
- **L3** leaves (`runtime/`, `workspace/`, `cli-plan/`, `polish/docs/`):
  `Op=runtime_multi|runtime_once`, `testbin.Ensure`, `label: heavy`.
- Nested `runsuite/` is its own L2 root (own Request/Run).
- All fixtures and store roots use `t.TempDir()` — never write into the user cache.
- Product default store path: `$CacheHome/doctest/leaf-cache/v1`; L3 tests isolate via
  `DOCTEST_CACHE_HOME` and/or `DOCTEST_LEAF_CACHE`.

## Steps

1. L2 library: leaf/group `Setup` builds fixtures; Run calls leafcache APIs only.
2. L3 runtime: `Setup` builds doctest binary + mini fixture; nested `doctest test` ×N.
3. L2 runsuite: twin trees + store seed; FormatLeafIdentity / PreparePassPlan / RecordPasses.
4. L3 workspace: multi-tree module; nested `doctest test <mod>/...`.
5. L3 cli-plan: multi-arg `doctest test tree-a tree-b`.
6. Assert keys, store hits, or summary Cached (L3 only).

## Context

### Expected public API surface (sealed product)

- `const AlgoVersion = "v1"`
- `type KeyInput struct { ModuleRoot, TreeRoot, LeafDir, GoVersion string }`
- `func ComputeLeafKey(in KeyInput) (string, error)`
- `func NewStore(root string) (*Store, error)`
- `func (s *Store) GetPass(key string) (bool, error)` — false when missing
- `func (s *Store) PutPass(key string) error` — explicit only; never auto on fail

### Spine hashing (must be reflected by ComputeLeafKey)

| Layer | Included |
|-------|----------|
| Leaf | `SETUP.md` Go + `ASSERT.md` Go |
| Ancestors | each `SETUP.md` Go up to tree root |
| Tree root | `DOCTEST.md` Go block |
| Module | `go.mod`, `go.sum` if present |
| Local pkgs | import closure under ModuleRoot reachable from spine |
| Local replace | replace target `go.mod` + local sources when in closure |
| Remote | **not** file-hashed; only go.mod/go.sum identity |
| Context | `GoVersion` string + algo version |

### Fixture layout helpers produce

**base** (`Flavor=base`):

```
$WorkDir/
  app/
    go.mod                  # module example.com/app
    pkg/helper/helper.go    # imported by leaf assert
    unrelated/noise.go      # NOT imported
    tests/feature/          # TreeRoot
      DOCTEST.md
      SETUP.md
      group/SETUP.md
      group/leaf/SETUP.md + ASSERT.md
```

**replace** (`Flavor=replace`): base + sibling `lib/` module and
`replace example.com/lib => ../lib` in app `go.mod`; leaf assert imports
`example.com/lib`.

**remote** (`Flavor=remote`): base + a remote-looking require line in go.mod
and a non-local fake tree `remote-src/example.com/remote@v1.0.0/*.go` that is
**not** a replace target — mutating those files must not change the key.

### Mutations (`applyMutation`)

| Name | Effect |
|------|--------|
| `leaf_assert` | append a comment to leaf ASSERT Go |
| `ancestor_setup` | append a comment to group SETUP Go |
| `local_imported` | change `pkg/helper/helper.go` body |
| `local_unrelated` | change `unrelated/noise.go` body |
| `replace_lib_src` | change `lib/lib.go` body |
| `replace_lib_gomod` | change `lib/go.mod` (e.g. go version line) |
| `remote_proxy_file` | change file under `remote-src/...` |

### P2 runtime helpers (see `runtime/` branch)

- `preparePassFixture` / `prepareFailFixture` — mini trees via `testtree`.
- `isolateRuntimeEnv` — isolated `DOCTEST_CACHE_HOME` / `DOCTEST_LEAF_CACHE`.
  `Run` adds a **fresh GOCACHE per invocation** so go testcache cannot inflate Cached.
- CLI flags under test: `-count`, `-a`, `--no-leaf-cache`.

### Polish helpers (see `polish/` + `key/tree-identity/`)

- `prepareTwoSiblingPassLeaves`, `prepareLocalDepPassFixture`, `prepareTwinTrees`
- `preparePartialPackageDepsFixture` — 4 pkgs + 3 leaves (shared a/b/c vs alone/d)
- `applyPolishMutation` / `MutateAfterRun` for mid-sequence edits
- `compute_two_inputs`, `runtime_once` ops

### RunSuite multi-prep (nested tree `runsuite/`)

- Nested `DOCTEST.md` root — see `tests/leaf-cache/runsuite/` for Request/Run.
- Implementer API: `FormatLeafIdentity`, `LeafRef`, `PassPlan`, `PreparePassPlan`, `RecordPasses`
- Parent `prepareTwinTrees` is mirrored inside the nested tree setup.

### Workspace multi-tree product (branch `workspace/` — P2)

- Helpers: `prepareWorkspaceTwoTrees`, `prepareWorkspaceAllPass`,
  `prepareWorkspacePartialFail`, `prepareWorkspaceSameRelpathPassFail`,
  `mustWorkspacePattern` / `workspaceDotDotDot` — under `workspace/SETUP.md`.
- Nested CLI uses parent `runtime_multi` + isolated leaf-cache env + fresh GOCACHE.
- Product path under test: `doctest test <mod>/...` → `RunWorkspace` /
  `finishWorkspaceGoTest` (not library-only PreparePassPlan).

### CLI one plan multi-arg (branch `cli-plan/` — P3)

- Helpers under `cli-plan/SETUP.md`: `prepareMultiArgAllPass` /
  `prepareMultiArgTwoTrees` / `multiArgTwoTrees` (same two-tree module layout
  as workspace; local copies because sibling `workspace/` helpers are not
  inherited across branches).
- Product path under test: `doctest test tree-a tree-b` — leaf-cache warm +
  `-count` bypass with the same store isolation as runtime/workspace.
- Asserts use `sumCachedCount` so N× per-tree summaries and a future single
  aggregated summary both count as total warm skips.

### Stream PutPass + grey progress dots (branch `runtime/stream-pass`, `runtime/progress-dots`)

- `prepareStreamInterruptFixture` — two instant-pass leaves + one hang leaf;
  `InterruptAfterDots` + SIGINT on first `runtime_multi` invocation; mutation
  `stream_unhang` rewrites hang ASSERT so run2 can finish.
- `preparePassFailMixFixture` — N pass + M fail leaves for red fail-dot lock.
- Grey/red progress-dot helpers live in root `DOCTEST.md` (`countGrayProgressDots`,
  `countRedProgressDots`, `progressPrefix`).

Still out of scope: full-repo CI stress, remote file hashing, exhaustive
hub-only fixtures (hub shares finishWorkspaceGoTest).

```go
import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/testtree"
)

// mdFence is built without a literal triple-backtick sequence in this source
// file so the SETUP.md markdown parser sees only one Go fence.
func mdFence() string {
	return strings.Repeat("`", 3)
}

func goFenceOpen() string {
	return mdFence() + "go\n"
}

func fenceClose() string {
	return mdFence() + "\n"
}

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	if req.GoVersion == "" {
		req.GoVersion = "go1.25.0"
	}
	return nil
}

// ensureWorkspace creates WorkDir and builds the fixture for req.Flavor.
func ensureWorkspace(t *testing.T, req *Request) error {
	t.Helper()
	if req.WorkDir == "" {
		req.WorkDir = t.TempDir()
	}
	flavor := req.Flavor
	if flavor == "" {
		flavor = "base"
	}
	switch flavor {
	case "base":
		return writeBaseFixture(t, req, false, false)
	case "replace":
		return writeBaseFixture(t, req, true, false)
	case "remote":
		return writeBaseFixture(t, req, false, true)
	default:
		return fmt.Errorf("unknown Flavor %q", flavor)
	}
}

func writeBaseFixture(t *testing.T, req *Request, withReplace, withRemote bool) error {
	t.Helper()
	app := filepath.Join(req.WorkDir, "app")
	lib := filepath.Join(req.WorkDir, "lib")
	tree := filepath.Join(app, "tests", "feature")
	leaf := filepath.Join(tree, "group", "leaf")

	if err := os.MkdirAll(filepath.Join(app, "pkg", "helper"), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(app, "unrelated"), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(leaf, 0o755); err != nil {
		return err
	}

	goMod := "module example.com/app\n\ngo 1.22\n"
	if withReplace {
		goMod += "\nrequire example.com/lib v0.0.0\n\nreplace example.com/lib => ../lib\n"
	}
	if withRemote {
		// Remote require identity only — no replace, no vendored source.
		goMod += "\nrequire example.com/remote v1.0.0\n"
	}
	if err := os.WriteFile(filepath.Join(app, "go.mod"), []byte(goMod), 0o644); err != nil {
		return err
	}
	// Minimal go.sum placeholder so presence is part of the module identity.
	if err := os.WriteFile(filepath.Join(app, "go.sum"), []byte("example.com/remote v1.0.0 h1:placeholder\n"), 0o644); err != nil {
		return err
	}

	helperSrc := "package helper\n\n// Answer is imported by the leaf assert spine.\nfunc Answer() int { return 42 }\n"
	if err := os.WriteFile(filepath.Join(app, "pkg", "helper", "helper.go"), []byte(helperSrc), 0o644); err != nil {
		return err
	}
	noiseSrc := "package unrelated\n\nfunc Noise() string { return \"noise\" }\n"
	if err := os.WriteFile(filepath.Join(app, "unrelated", "noise.go"), []byte(noiseSrc), 0o644); err != nil {
		return err
	}

	if withReplace {
		if err := os.MkdirAll(lib, 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(lib, "go.mod"), []byte("module example.com/lib\n\ngo 1.22\n"), 0o644); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(lib, "lib.go"), []byte("package lib\n\nfunc Val() int { return 7 }\n"), 0o644); err != nil {
			return err
		}
	}

	if withRemote {
		remoteDir := filepath.Join(req.WorkDir, "remote-src", "example.com", "remote@v1.0.0")
		if err := os.MkdirAll(remoteDir, 0o755); err != nil {
			return err
		}
		// Not reachable via replace or local import path — must not enter the DAG.
		if err := os.WriteFile(filepath.Join(remoteDir, "remote.go"), []byte("package remote\n\nfunc X() int { return 1 }\n"), 0o644); err != nil {
			return err
		}
	}

	// Root DOCTEST.md — only the final Go block is spine-relevant.
	doctestMD := "# Feature\n\n## Version\n0.0.2\n\n# DSN (Domain Specific Notion)\n\n" +
		"Fixture tree for leaf-cache key hashing.\n\n" +
		goFenceOpen() +
		"import \"testing\"\n\n" +
		"type Request struct{ N int }\n" +
		"type Response struct{ Ok bool }\n\n" +
		"func Run(t *testing.T, req *Request) (*Response, error) {\n" +
		"\treturn &Response{Ok: true}, nil\n" +
		"}\n" +
		fenceClose()
	if err := os.WriteFile(filepath.Join(tree, "DOCTEST.md"), []byte(doctestMD), 0o644); err != nil {
		return err
	}

	rootSetup := "# Scenario\n\n**Feature**: fixture root\n\n" +
		mdFence() + "\nfixture root setup\n" + fenceClose() + "\n" +
		"## Steps\n1. Root setup.\n\n" +
		goFenceOpen() +
		"import \"testing\"\n\n" +
		"func Setup(t *testing.T, req *Request) error {\n" +
		"\treq.N = 1\n" +
		"\treturn nil\n" +
		"}\n" +
		fenceClose()
	if err := os.WriteFile(filepath.Join(tree, "SETUP.md"), []byte(rootSetup), 0o644); err != nil {
		return err
	}

	groupSetup := "# Scenario\n\n**Feature**: fixture group\n\n" +
		mdFence() + "\ngroup setup\n" + fenceClose() + "\n" +
		"## Steps\n1. Group setup.\n\n" +
		goFenceOpen() +
		"import \"testing\"\n\n" +
		"func Setup(t *testing.T, req *Request) error {\n" +
		"\treq.N = 2\n" +
		"\treturn nil\n" +
		"}\n" +
		fenceClose()
	if err := os.WriteFile(filepath.Join(tree, "group", "SETUP.md"), []byte(groupSetup), 0o644); err != nil {
		return err
	}

	leafSetup := "# Scenario\n\n**Feature**: fixture leaf\n\n" +
		mdFence() + "\nleaf setup\n" + fenceClose() + "\n" +
		"## Steps\n1. Leaf setup.\n\n" +
		goFenceOpen() +
		"import \"testing\"\n\n" +
		"func Setup(t *testing.T, req *Request) error {\n" +
		"\treq.N = 3\n" +
		"\treturn nil\n" +
		"}\n" +
		fenceClose()
	if err := os.WriteFile(filepath.Join(leaf, "SETUP.md"), []byte(leafSetup), 0o644); err != nil {
		return err
	}

	assertImports := "\t\"example.com/app/pkg/helper\"\n"
	assertBody := "\tv := helper.Answer()\n\tif v != 42 {\n\t\tt.Fatalf(\"answer = %d\", v)\n\t}\n\t_ = req\n\t_ = resp\n\t_ = err\n"
	if withReplace {
		assertImports = "\t\"example.com/app/pkg/helper\"\n\t\"example.com/lib\"\n"
		assertBody = "\tv := helper.Answer() + lib.Val()\n\tif v != 49 {\n\t\tt.Fatalf(\"answer = %d\", v)\n\t}\n\t_ = req\n\t_ = resp\n\t_ = err\n"
	}

	assertMD := "## Expected\n\n- helper answer is 42\n\n" +
		goFenceOpen() +
		"import (\n\t\"testing\"\n" + assertImports + ")\n\n" +
		"func Assert(t *testing.T, req *Request, resp *Response, err error) {\n" +
		assertBody +
		"}\n" +
		fenceClose()
	if err := os.WriteFile(filepath.Join(leaf, "ASSERT.md"), []byte(assertMD), 0o644); err != nil {
		return err
	}

	req.ModuleRoot = app
	req.TreeRoot = tree
	req.LeafDir = leaf
	return nil
}

func applyMutation(t *testing.T, req *Request) error {
	t.Helper()
	switch req.Mutation {
	case "leaf_assert":
		return mutateFileContent(filepath.Join(req.LeafDir, "ASSERT.md"), "v := helper.Answer()", "v := helper.Answer() // mutated")
	case "ancestor_setup":
		return mutateFileContent(filepath.Join(req.TreeRoot, "group", "SETUP.md"), "req.N = 2", "req.N = 2 // mutated")
	case "local_imported":
		p := filepath.Join(req.ModuleRoot, "pkg", "helper", "helper.go")
		return os.WriteFile(p, []byte("package helper\n\nfunc Answer() int { return 99 }\n"), 0o644)
	case "local_unrelated":
		p := filepath.Join(req.ModuleRoot, "unrelated", "noise.go")
		return os.WriteFile(p, []byte("package unrelated\n\nfunc Noise() string { return \"changed\" }\n"), 0o644)
	case "replace_lib_src":
		p := filepath.Join(req.WorkDir, "lib", "lib.go")
		return os.WriteFile(p, []byte("package lib\n\nfunc Val() int { return 100 }\n"), 0o644)
	case "replace_lib_gomod":
		p := filepath.Join(req.WorkDir, "lib", "go.mod")
		return os.WriteFile(p, []byte("module example.com/lib\n\ngo 1.23\n"), 0o644)
	case "remote_proxy_file":
		p := filepath.Join(req.WorkDir, "remote-src", "example.com", "remote@v1.0.0", "remote.go")
		return os.WriteFile(p, []byte("package remote\n\nfunc X() int { return 999 }\n"), 0o644)
	case "":
		return fmt.Errorf("Mutation is empty")
	default:
		return fmt.Errorf("unknown Mutation %q", req.Mutation)
	}
}

func mutateFileContent(path, old, new string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	s := string(b)
	// Prefer in-body marker; fall back to replace-flavor assert form.
	if strings.Contains(s, old) {
		s = strings.Replace(s, old, new, 1)
		return os.WriteFile(path, []byte(s), 0o644)
	}
	if old == "v := helper.Answer()" {
		altOld := "v := helper.Answer() + lib.Val()"
		altNew := "v := helper.Answer() + lib.Val() // mutated"
		if strings.Contains(s, altOld) {
			s = strings.Replace(s, altOld, altNew, 1)
			return os.WriteFile(path, []byte(s), 0o644)
		}
	}
	return fmt.Errorf("mutation marker %q not found in %s", old, path)
}

// isolateRuntimeEnv returns env that isolates leaf-cache store under a temp
// CacheHome. GOCACHE is intentionally omitted here; runtime_multi injects a
// fresh GOCACHE per invocation so go testcache cannot produce cross-run Cached.
func isolateRuntimeEnv(t *testing.T) []string {
	t.Helper()
	cacheHome := t.TempDir()
	leafRoot := filepath.Join(cacheHome, "doctest", "leaf-cache", "v1")
	return []string{
		"DOCTEST_CACHE_HOME=" + cacheHome,
		"DOCTEST_LEAF_CACHE=" + leafRoot,
		// Stable session so nested gen paths do not thrash across the two runs.
		"DOCTEST_SESSION_ID=leaf-cache-runtime-stable",
	}
}

// preparePassFixture writes a mini always-pass doctest tree with passCount leaves.
func preparePassFixture(t *testing.T, passCount int) string {
	t.Helper()
	dir := t.TempDir()
	if passCount <= 0 {
		passCount = 1
	}
	testtree.WritePassFailTree(t, dir, passCount, 0)
	return dir
}

// prepareFailFixture writes a mini always-fail doctest tree with failCount leaves.
func prepareFailFixture(t *testing.T, failCount int) string {
	t.Helper()
	dir := t.TempDir()
	if failCount <= 0 {
		failCount = 1
	}
	testtree.WritePassFailTree(t, dir, 0, failCount)
	return dir
}

// prepareStreamInterruptFixture writes one instant-pass leaf and one hang leaf
// (long sleep). SIGINT after the pass progress dot cannot reach end-of-run
// recordLeafCachePasses while hang is still running.
func prepareStreamInterruptFixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	testtree.WriteMinimalRunnableTree(t, dir, []testtree.LeafSpec{
		{
			Name: "leaf_a",
			AssertGo: `import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	// stream_leaf_a
	_ = req
	_ = resp
	_ = err
}`,
		},
		{
			Name: "leaf_hang",
			AssertGo: `import (
	"testing"
	"time"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	// stream_hang_marker
	_ = req
	_ = resp
	_ = err
	time.Sleep(10 * time.Minute)
}`,
		},
	})
	return dir
}

// preparePassFailMixFixture writes passCount always-pass + failCount always-fail leaves.
func preparePassFailMixFixture(t *testing.T, passCount, failCount int) string {
	t.Helper()
	dir := t.TempDir()
	if passCount < 0 {
		passCount = 0
	}
	if failCount < 0 {
		failCount = 0
	}
	if passCount+failCount == 0 {
		passCount = 1
	}
	testtree.WritePassFailTree(t, dir, passCount, failCount)
	return dir
}

// prepareTwoSiblingPassLeaves writes a tree with leaf_a and leaf_b (always pass).
// Markers in ASSERT Go allow polish_edit_leaf_a mutation.
func prepareTwoSiblingPassLeaves(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	testtree.WriteMinimalRunnableTree(t, dir, []testtree.LeafSpec{
		{
			Name: "leaf_a",
			AssertGo: `import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	// marker_leaf_a
	_ = req
	_ = resp
	_ = err
}`,
		},
		{
			Name: "leaf_b",
			AssertGo: `import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	// marker_leaf_b
	_ = req
	_ = resp
	_ = err
}`,
		},
	})
	return dir
}

// prepareLocalDepPassFixture writes a mini module + single pass leaf that imports
// example.com/app/pkg/helper. FixtureDir is the doctest tree under app/tests/feature.
func prepareLocalDepPassFixture(t *testing.T, req *Request) string {
	t.Helper()
	work := t.TempDir()
	app := filepath.Join(work, "app")
	tree := filepath.Join(app, "tests", "feature")
	leaf := filepath.Join(tree, "leaf")
	if err := os.MkdirAll(filepath.Join(app, "pkg", "helper"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(leaf, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(app, "go.mod"), []byte("module example.com/app\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(app, "pkg", "helper", "helper.go"), []byte("package helper\n\nfunc Answer() int { return 42 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Minimal runnable tree files
	doctestBody := testtree.MinimalDOCTEST(testtree.MinimalRunGo())
	if err := os.WriteFile(filepath.Join(tree, "DOCTEST.md"), []byte(doctestBody), 0o644); err != nil {
		t.Fatal(err)
	}
	setupBody := "## Steps\n1. leaf\n\n" + goFenceOpen() + "import \"testing\"\n\nfunc Setup(t *testing.T, req *Request) error { return nil }\n" + fenceClose()
	if err := os.WriteFile(filepath.Join(leaf, "SETUP.md"), []byte(setupBody), 0o644); err != nil {
		t.Fatal(err)
	}
	assertBody := "## Expected\n- helper\n\n" + goFenceOpen() +
		"import (\n\t\"testing\"\n\t\"example.com/app/pkg/helper\"\n)\n\n" +
		"func Assert(t *testing.T, req *Request, resp *Response, err error) {\n" +
		"\tif helper.Answer() != 42 {\n\t\tt.Fatalf(\"answer\")\n\t}\n" +
		"\t_ = req\n\t_ = resp\n\t_ = err\n}\n" + fenceClose()
	if err := os.WriteFile(filepath.Join(leaf, "ASSERT.md"), []byte(assertBody), 0o644); err != nil {
		t.Fatal(err)
	}
	req.WorkDir = work
	req.ModuleRoot = app
	req.TreeRoot = tree
	req.LeafDir = leaf
	req.FixtureDir = tree
	return tree
}

// prepareTwinTrees writes two content-identical single-leaf trees under work.
// Same relative path "leaf"; different absolute TreeRoots.
func prepareTwinTrees(t *testing.T, req *Request) {
	t.Helper()
	work := t.TempDir()
	writeOneLeafTree := func(root string) {
		testtree.WriteMinimalRunnableTree(t, root, []testtree.LeafSpec{
			{
				Name: "leaf",
				AssertGo: `import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	// twin_marker
	_ = req
	_ = resp
	_ = err
}`,
			},
		})
	}
	treeA := filepath.Join(work, "treeA")
	treeB := filepath.Join(work, "treeB")
	if err := os.MkdirAll(treeA, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(treeB, 0o755); err != nil {
		t.Fatal(err)
	}
	writeOneLeafTree(treeA)
	writeOneLeafTree(treeB)
	req.WorkDir = work
	req.FixtureDir = treeA
	req.FixtureB = treeB
	req.TreeRoot = treeA
	req.LeafDir = filepath.Join(treeA, "leaf")
	req.ModuleRoot = treeA // no go.mod; KeyForLeaf uses tree root
	req.TreeRootB = treeB
	req.LeafDirB = filepath.Join(treeB, "leaf")
	req.ModuleRootB = treeB
}

// preparePartialPackageDepsFixture builds 4 packages + 3 leaves:
//   shared/{a,b,c} used by leaf-ab-1 and leaf-ab-2; alone/d used only by leaf-d.
// TreeRoot is app/tests/feature. Leaves call Version() so packages are in the key DAG.
func preparePartialPackageDepsFixture(t *testing.T, req *Request) string {
	t.Helper()
	work := t.TempDir()
	app := filepath.Join(work, "app")
	tree := filepath.Join(app, "tests", "feature")
	writePkg := func(rel, pkg, version string) {
		dir := filepath.Join(app, filepath.FromSlash(rel))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		body := fmt.Sprintf("package %s\n\nfunc Version() string { return %q }\n", pkg, version)
		if err := os.WriteFile(filepath.Join(dir, pkg+".go"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	writePkg("shared/a", "a", "A")
	writePkg("shared/b", "b", "B")
	writePkg("shared/c", "c", "C")
	writePkg("alone/d", "d", "D")
	if err := os.WriteFile(filepath.Join(app, "go.mod"), []byte("module example.com/partialcache\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(tree, 0o755); err != nil {
		t.Fatal(err)
	}
	doctestBody := testtree.MinimalDOCTEST(testtree.MinimalRunGo())
	if err := os.WriteFile(filepath.Join(tree, "DOCTEST.md"), []byte(doctestBody), 0o644); err != nil {
		t.Fatal(err)
	}
	writeLeaf := func(name, assertGo string) {
		leaf := filepath.Join(tree, name)
		if err := os.MkdirAll(leaf, 0o755); err != nil {
			t.Fatal(err)
		}
		setupBody := "## Steps\n1. leaf\n\n" + goFenceOpen() +
			"import \"testing\"\n\nfunc Setup(t *testing.T, req *Request) error { return nil }\n" + fenceClose()
		if err := os.WriteFile(filepath.Join(leaf, "SETUP.md"), []byte(setupBody), 0o644); err != nil {
			t.Fatal(err)
		}
		assertBody := "## Expected\n- versions\n\n" + goFenceOpen() + assertGo + fenceClose()
		if err := os.WriteFile(filepath.Join(leaf, "ASSERT.md"), []byte(assertBody), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// A/B share a,b,c; C uses alone/d only. Non-empty Version checks keep imports live.
	abAssert := `import (
	"testing"
	"example.com/partialcache/shared/a"
	"example.com/partialcache/shared/b"
	"example.com/partialcache/shared/c"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if a.Version() == "" || b.Version() == "" || c.Version() == "" {
		t.Fatal("empty shared version")
	}
	_ = req
	_ = resp
	_ = err
}
`
	// Alias alone package: bare "d" collides with injected d *session.Doctest.
	dAssert := `import (
	"testing"
	aloned "example.com/partialcache/alone/d"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if aloned.Version() == "" {
		t.Fatal("empty alone version")
	}
	_ = req
	_ = resp
	_ = err
}
`
	writeLeaf("leaf-ab-1", abAssert)
	writeLeaf("leaf-ab-2", abAssert)
	writeLeaf("leaf-d", dAssert)
	req.WorkDir = work
	req.ModuleRoot = app
	req.TreeRoot = tree
	req.FixtureDir = tree
	return tree
}

// applyPolishMutation applies P3 runtime mid-run mutations.
func applyPolishMutation(t *testing.T, req *Request) error {
	t.Helper()
	switch req.Mutation {
	case "polish_edit_leaf_a":
		p := filepath.Join(req.FixtureDir, "leaf_a", "ASSERT.md")
		return mutateFileContent(p, "// marker_leaf_a", "// marker_leaf_a mutated")
	case "polish_edit_local_dep":
		p := filepath.Join(req.ModuleRoot, "pkg", "helper", "helper.go")
		return os.WriteFile(p, []byte("package helper\n\nfunc Answer() int { return 99 }\n"), 0o644)
	case "polish_edit_alone_d":
		p := filepath.Join(req.ModuleRoot, "alone", "d", "d.go")
		return os.WriteFile(p, []byte("package d\n\nfunc Version() string { return \"D2\" }\n"), 0o644)
	case "polish_edit_shared_a":
		p := filepath.Join(req.ModuleRoot, "shared", "a", "a.go")
		return os.WriteFile(p, []byte("package a\n\nfunc Version() string { return \"A2\" }\n"), 0o644)
	case "stream_unhang":
		// Replace hang sleep ASSERT with an instant pass so run2 can finish.
		// Content change also ensures hang leaf is not treated as a warm hit
		// from a non-existent hang key (hang never PutPass'd).
		p := filepath.Join(req.FixtureDir, "leaf_hang", "ASSERT.md")
		body := "## Expected\n- unhung pass\n\n" + goFenceOpen() +
			"import \"testing\"\n\n" +
			"func Assert(t *testing.T, req *Request, resp *Response, err error) {\n" +
			"\t// stream_unhang_pass\n" +
			"\t_ = req\n\t_ = resp\n\t_ = err\n}\n" + fenceClose()
		return os.WriteFile(p, []byte(body), 0o644)
	case "":
		return fmt.Errorf("polish Mutation is empty")
	default:
		// Fall back to P1 applyMutation names if used mid-run.
		return applyMutation(t, req)
	}
}
```
