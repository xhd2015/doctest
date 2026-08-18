package build

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
)

// Path-scope contract (both phases):
//
//  1. Generate: only touch gen content under the selected subpath (with or
//     without ...). Content outside that path must not be created or rewritten.
//  2. go test: only execute packages/cases under that same subpath. Mid vs
//     sibling scopes must not collapse to the same go-test plan.
//
// These tests lock the intent before the layout/executor rework. Prefer RED
// on shared-root suite plans and rewrites outside the mid prefix.

func pathScopeMidSiblingTree(t *testing.T) (mod, tree, mid, sibling string) {
	t.Helper()
	mod = t.TempDir()
	if err := os.WriteFile(filepath.Join(mod, "go.mod"), []byte("module pathscope\ngo 1.21\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	tree = filepath.Join(mod, "tree")
	writeRootHarness(t, tree, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) { return &Response{}, nil }
`, `
func Setup(t *testing.T, d *session.Doctest, req *Request) error { return nil }
`)
	mid = filepath.Join(tree, "mid")
	sibling = filepath.Join(tree, "sibling")
	writeTreeFile(t, mid, "two/SETUP.md", setupDoc(`
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Logf("MARKER:MID_LEAF")
	return nil
}
`))
	writeTreeFile(t, mid, "two/ASSERT.md", assertDoc(`
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}
`))
	writeTreeFile(t, sibling, "one/SETUP.md", setupDoc(`
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Logf("MARKER:SIBLING_LEAF")
	return nil
}
`))
	writeTreeFile(t, sibling, "one/ASSERT.md", assertDoc(`
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}
`))
	return mod, tree, mid, sibling
}

func pathScopeHashTree(t *testing.T, root string) map[string]string {
	t.Helper()
	out := make(map[string]string)
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(b)
		out[rel] = hex.EncodeToString(sum[:])
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// bookkeepingRel is gen-root metadata that may change every run even when path
// scope is correct (manifest, tidy stamp, module file).
func pathScopeIsBookkeeping(rel string) bool {
	switch rel {
	case "go.mod", "go.sum", "doctest.gen-manifest", "doctest.tidy-done", core.ExposeMaterializedList:
		return true
	default:
		return false
	}
}

// underSelectedPrefix reports whether a gen-relative path is content for the
// selected source subpath (e.g. tree/mid/...).
func pathScopeUnderSelected(genRel, treeRel, selectedUnderTree string) bool {
	genRel = filepath.ToSlash(genRel)
	prefix := filepath.ToSlash(filepath.Join(treeRel, selectedUnderTree))
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	return genRel == strings.TrimSuffix(prefix, "/") || strings.HasPrefix(genRel, prefix)
}

func pathScopeExtractGoTestPlans(out string) []string {
	var plans []string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "&& go test") || strings.HasPrefix(line, "go test ") {
			plans = append(plans, line)
		}
	}
	return plans
}

// --- Phase 1: generate ---

// TestPathScope_Generate_ColdMidDoesNotEmitSibling: first gen for mid/... must
// not create packages for sibling under the gen root.
func TestPathScope_Generate_ColdMidDoesNotEmitSibling(t *testing.T) {
	_, tree, mid, _ := pathScopeMidSiblingTree(t)
	gen := filepath.Join(t.TempDir(), "gen")
	var stderr bytes.Buffer
	if err := Test(tree, core.Options{
		GenDir:       gen,
		SubDir:       mid,
		BypassGoTest: true,
		Stderr:       &stderr,
	}); err != nil {
		t.Fatalf("mid gen: %v\n%s", err, stderr.String())
	}
	// Sibling leaf packages must not appear.
	err := filepath.Walk(gen, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(gen, p)
		rel = filepath.ToSlash(rel)
		if strings.Contains(rel, "/sibling/") || strings.HasPrefix(rel, "sibling/") || strings.Contains(rel, "tree/sibling") {
			t.Errorf("cold mid gen wrote out-of-scope sibling path: %s", rel)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	// Mid leaf should exist.
	midLeaf := filepath.Join(gen, "tree", "mid", "two", "leaf.go")
	if _, err := os.Stat(midLeaf); err != nil {
		t.Fatalf("expected mid leaf under gen: %v (stderr:\n%s)", err, stderr.String())
	}
}

// TestPathScope_Generate_FullThenMidDoesNotRewriteOutsideMid: after full-tree
// gen, a scoped mid gen must leave every non-bookkeeping file outside tree/mid/
// byte-identical (sibling leaves, tree suite/registry, workspace fan-in, …).
func TestPathScope_Generate_FullThenMidDoesNotRewriteOutsideMid(t *testing.T) {
	_, tree, mid, _ := pathScopeMidSiblingTree(t)
	gen := filepath.Join(t.TempDir(), "gen")
	var stderr bytes.Buffer
	if err := Test(tree, core.Options{
		GenDir:       gen,
		BypassGoTest: true,
		Stderr:       &stderr,
	}); err != nil {
		t.Fatalf("full gen: %v\n%s", err, stderr.String())
	}
	before := pathScopeHashTree(t, gen)

	stderr.Reset()
	if err := Test(tree, core.Options{
		GenDir:       gen,
		SubDir:       mid,
		BypassGoTest: true,
		Stderr:       &stderr,
	}); err != nil {
		t.Fatalf("mid gen: %v\n%s", err, stderr.String())
	}
	after := pathScopeHashTree(t, gen)

	var rewritten []string
	for rel, sum := range before {
		if pathScopeIsBookkeeping(rel) {
			continue
		}
		if pathScopeUnderSelected(rel, "tree", "mid") {
			continue // in-scope for mid may be rewritten
		}
		if after[rel] != sum {
			rewritten = append(rewritten, rel)
		}
	}
	// Files only present after mid gen outside mid are also out-of-scope writes.
	for rel := range after {
		if pathScopeIsBookkeeping(rel) || pathScopeUnderSelected(rel, "tree", "mid") {
			continue
		}
		if _, ok := before[rel]; !ok {
			rewritten = append(rewritten, rel+" (new outside mid)")
		}
	}
	if len(rewritten) > 0 {
		t.Fatalf("mid gen rewrote or created content outside tree/mid (must not touch out-of-scope gen):\n  %s\nstderr:\n%s",
			strings.Join(rewritten, "\n  "), stderr.String())
	}
}

// TestPathScope_Generate_FullThenMidSiblingLeafUnchanged: weaker leaf-only
// check (sibling leaf.go content stable) — still required even if shared suite
// layout changes later.
func TestPathScope_Generate_FullThenMidSiblingLeafUnchanged(t *testing.T) {
	_, tree, mid, _ := pathScopeMidSiblingTree(t)
	gen := filepath.Join(t.TempDir(), "gen")
	if err := Test(tree, core.Options{GenDir: gen, BypassGoTest: true}); err != nil {
		t.Fatal(err)
	}
	sibLeaf := filepath.Join(gen, "tree", "sibling", "one", "leaf.go")
	orig, err := os.ReadFile(sibLeaf)
	if err != nil {
		t.Fatal(err)
	}
	planted := append(append([]byte(nil), orig...), []byte("\n// PLANT_SIBLING_MARKER\n")...)
	if err := os.WriteFile(sibLeaf, planted, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Test(tree, core.Options{GenDir: gen, SubDir: mid, BypassGoTest: true}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(sibLeaf)
	if err != nil {
		t.Fatalf("sibling leaf.go missing after mid gen (out-of-scope must not be deleted): %v", err)
	}
	if !bytes.Equal(got, planted) {
		t.Fatalf("sibling leaf.go rewritten by mid gen (out-of-scope content must stay)\n--- before plant len %d after mid len %d", len(planted), len(got))
	}
}

// --- Phase 2: go test ---

// TestPathScope_GoTest_MidExcludesSiblingMarkers: running mid must execute mid
// leaf and must not execute sibling leaf.
func TestPathScope_GoTest_MidExcludesSiblingMarkers(t *testing.T) {
	_, tree, mid, _ := pathScopeMidSiblingTree(t)
	gen := filepath.Join(t.TempDir(), "gen")
	var stdout, stderr bytes.Buffer
	if err := Test(tree, core.Options{
		GenDir:  gen,
		SubDir:  mid,
		Verbose: true,
		Count:   1,
		Stdout:  &stdout,
		Stderr:  &stderr,
	}); err != nil {
		t.Fatalf("mid test: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	out := stdout.String() + "\n" + stderr.String()
	if !strings.Contains(out, "MARKER:MID_LEAF") {
		t.Fatalf("want MID leaf executed\n%s", out)
	}
	if strings.Contains(out, "MARKER:SIBLING_LEAF") {
		t.Fatalf("sibling must not run under mid scope\n%s", out)
	}
}

// TestPathScope_GoTest_MidVsSiblingDistinctPathScopedPlans: mid and sibling
// scopes must produce different go-test plans, each path-scoped (pattern under
// the selected branch), not the same shared root suite for both.
func TestPathScope_GoTest_MidVsSiblingDistinctPathScopedPlans(t *testing.T) {
	_, tree, mid, sibling := pathScopeMidSiblingTree(t)
	gen := filepath.Join(t.TempDir(), "gen")

	runScoped := func(sub string) (plans []string, out string) {
		t.Helper()
		var stdout, stderr bytes.Buffer
		err := Test(tree, core.Options{
			GenDir:  gen,
			SubDir:  sub,
			Verbose: true,
			Count:   1,
			Stdout:  &stdout,
			Stderr:  &stderr,
		})
		out = stdout.String() + "\n" + stderr.String()
		if err != nil {
			t.Fatalf("test sub=%s: %v\n%s", sub, err, out)
		}
		return pathScopeExtractGoTestPlans(out), out
	}

	midPlans, midOut := runScoped(mid)
	sibPlans, sibOut := runScoped(sibling)

	if len(midPlans) == 0 {
		t.Fatalf("no go test plan line for mid:\n%s", midOut)
	}
	if len(sibPlans) == 0 {
		t.Fatalf("no go test plan line for sibling:\n%s", sibOut)
	}
	// Distinct plans: mid vs sibling must not collapse to the same command.
	if midPlans[0] == sibPlans[0] {
		t.Fatalf("mid and sibling scopes share the same go test plan (filter lost at go-test level):\n  mid: %s\n  sib: %s\n", midPlans[0], sibPlans[0])
	}
	// Path-scoped: ./…/mid/... and ./…/sibling/... (not hard-coded */suite).
	if !strings.Contains(midPlans[0], "mid") || !strings.Contains(midPlans[0], "/...") {
		t.Fatalf("mid go test plan want path ... under mid, got:\n  %s\nfull:\n%s", midPlans[0], midOut)
	}
	if !strings.Contains(sibPlans[0], "sibling") || !strings.Contains(sibPlans[0], "/...") {
		t.Fatalf("sibling go test plan want path ... under sibling, got:\n  %s\nfull:\n%s", sibPlans[0], sibOut)
	}
	if strings.Contains(midPlans[0], "/suite") && !strings.Contains(midPlans[0], "/...") {
		t.Fatalf("mid plan must not be */suite workaround: %s", midPlans[0])
	}
	if strings.Contains(midPlans[0], "__workspace/suite") {
		t.Fatalf("mid plan uses root workspace suite without mid path scope:\n  %s", midPlans[0])
	}
}
