# Scenario

**Feature**: L2 policy selection with fixture trees + synthetic changed path lists

```
# ephemeral fixture tree (no real git required)
write DOCTEST/SETUP/ASSERT leaves
  -> req.TreeDir + req.GitRoot + req.ChangedFiles
  -> core.FilterByChangedFiles | ChangedRunInfoForTree | ChangedDoctestMarkdownFiles

# map changed files to leaves
changed ASSERT.md -> one leaf | changed group SETUP.md -> descendant leaves
```

## Preconditions

- L2: fixture trees use valid `DOCTEST.md` with `Request`, `Response`, and `Run`.
- Changed path lists are relative to `GitRoot` (typically the temp parent of `tests/fixture`).
- No product binary; no `git init` for pure filter policy.

## Steps

1. Create an ephemeral fixture tree under `t.TempDir()`.
2. Set synthetic `ChangedFiles` paths for the leaf scenario.
3. `Run` applies core filter / vet-md APIs in-process.

## Context

- Helpers `createFlatTwoLeafTree`, `createSharedParentTwoLeafTree`, and path
  builders produce reproducible fixtures.
- `relUnderTree` builds changed paths relative to the synthetic git root.

```go
import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/testtree"
)

const tick = "\u0060"

const leafSetupGo = "import \"testing\"\n\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }"

const leafAssertGo = "import \"testing\"\n\nfunc Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}"

const fixtureSetupGo = "import \"testing\"\n\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }"

// policyFixture is a temp layout whose GitRoot is the synthetic repo parent and
// TreeDir is tests/fixture under it. No real git repository is required.
type policyFixture struct {
	RepoDir string // GitRoot
	TreeDir string // doctest root
}

func goFence(n int) string {
	return strings.Repeat(tick, n)
}

func goBlock(code string) string {
	fence := goFence(3)
	return fence + "go\n" + code + "\n" + fence + "\n"
}

func scenarioHeader(feature, snippet string) string {
	fence := goFence(3)
	return fmt.Sprintf("# Scenario\n\n**Feature**: %s\n\n%s\n%s\n%s\n\n", feature, fence, snippet, fence)
}

func writeLeaf(t *testing.T, root, name string) {
	t.Helper()
	leafDir := filepath.Join(root, name)
	if err := os.MkdirAll(leafDir, 0755); err != nil {
		t.Fatal(err)
	}
	setup := scenarioHeader("fixture leaf", "# leaf setup\nleaf -> run") + "## Steps\n1. setup\n\n" + goBlock(leafSetupGo)
	assert := "## Expected\n- passes\n\n" + goBlock(leafAssertGo)
	if err := os.WriteFile(filepath.Join(leafDir, "SETUP.md"), []byte(setup), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(leafDir, "ASSERT.md"), []byte(assert), 0644); err != nil {
		t.Fatal(err)
	}
}

func writeRootTree(t *testing.T, root string, withRootSetup bool) {
	t.Helper()
	if err := os.MkdirAll(root, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(testtree.MinimalRunGo())), 0644); err != nil {
		t.Fatal(err)
	}
	if withRootSetup {
		setup := scenarioHeader("fixture root", "# root setup\nroot -> shared context") + "## Steps\n1. root setup\n\n" + goBlock(fixtureSetupGo)
		if err := os.WriteFile(filepath.Join(root, "SETUP.md"), []byte(setup), 0644); err != nil {
			t.Fatal(err)
		}
	}
}

func createFlatTwoLeafTree(t *testing.T) policyFixture {
	t.Helper()
	repoDir := t.TempDir()
	treeDir := filepath.Join(repoDir, "tests", "fixture")
	writeRootTree(t, treeDir, true)
	writeLeaf(t, treeDir, "leaf_a")
	writeLeaf(t, treeDir, "leaf_b")
	return policyFixture{RepoDir: repoDir, TreeDir: treeDir}
}

func createSharedParentTwoLeafTree(t *testing.T) policyFixture {
	t.Helper()
	repoDir := t.TempDir()
	treeDir := filepath.Join(repoDir, "tests", "fixture")
	writeRootTree(t, treeDir, true)
	groupDir := filepath.Join(treeDir, "shared")
	if err := os.MkdirAll(groupDir, 0755); err != nil {
		t.Fatal(err)
	}
	groupSetup := scenarioHeader("shared parent", "# group setup\nshared -> both leaves") + "## Steps\n1. group setup\n\n" + goBlock(fixtureSetupGo)
	if err := os.WriteFile(filepath.Join(groupDir, "SETUP.md"), []byte(groupSetup), 0644); err != nil {
		t.Fatal(err)
	}
	writeLeaf(t, groupDir, "leaf_a")
	writeLeaf(t, groupDir, "leaf_b")
	return policyFixture{RepoDir: repoDir, TreeDir: treeDir}
}

// treeRel returns a changed-file path relative to the synthetic git root for a
// path under the fixture tree (e.g. leaf_a/ASSERT.md -> tests/fixture/leaf_a/ASSERT.md).
func treeRel(fx policyFixture, parts ...string) string {
	elems := append([]string{"tests", "fixture"}, parts...)
	return filepath.ToSlash(filepath.Join(elems...))
}

func applyPolicyBase(req *Request, fx policyFixture) {
	req.TreeDir = fx.TreeDir
	req.GitRoot = fx.RepoDir
	req.UseCLI = false
}

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.UseCLI = false
	return nil
}
```
