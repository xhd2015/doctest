package core

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestAssembleWorkspaceRegistryAndSuite(t *testing.T) {
	reg := AssembleWorkspaceRegistrySource()
	if !strings.Contains(reg, "type TreeEntry struct") {
		t.Fatalf("registry missing TreeEntry:\n%s", reg)
	}
	if !strings.Contains(reg, "func Register(e TreeEntry)") {
		t.Fatalf("registry missing Register:\n%s", reg)
	}
	if !strings.Contains(reg, "func All() []TreeEntry") {
		t.Fatalf("registry missing All:\n%s", reg)
	}
	if strings.Contains(reg, "IsolateProcessEnv") {
		t.Fatalf("registry must not include IsolateProcessEnv (process env mutation):\n%s", reg)
	}

	wreg := AssembleTreeWregSource(
		"libdoc/build/tests",
		false,
		"testcase/libdoc/build/tests/__registry",
		"testcase/libdoc/build/tests/__allleaves",
		WorkspaceRegistryImport(),
	)
	if !strings.Contains(wreg, "wsreg.Register") {
		t.Fatalf("wreg must register into workspace:\n%s", wreg)
	}
	if !strings.Contains(wreg, "treereg.All()") {
		t.Fatalf("wreg must iterate tree registry:\n%s", wreg)
	}
	if strings.Contains(wreg, "IsolateProcessEnv") || strings.Contains(wreg, "Setenv") {
		t.Fatalf("wreg must not mutate process env:\n%s", wreg)
	}
	if !strings.Contains(wreg, `Path:  "libdoc/build/tests"`) {
		t.Fatalf("wreg path:\n%s", wreg)
	}

	heavy := AssembleTreeWregSource("tests", true, "testcase/tests/__registry", "testcase/tests/__allleaves", WorkspaceRegistryImport())
	if !strings.Contains(heavy, "Heavy: true") {
		t.Fatalf("heavy wreg should set Heavy:\n%s", heavy)
	}

	suite := AssembleWorkspaceSuiteTestSource(WorkspaceRegistryImport(), WorkspaceAllTreesImport())
	if !strings.Contains(suite, "registry.All()") {
		t.Fatalf("suite must use workspace All:\n%s", suite)
	}
	// v1: serial trees (no t.Parallel) — nested selftests are not process-safe.
	if strings.Contains(suite, "t.Parallel()") {
		t.Fatalf("v1 suite must not Parallel trees:\n%s", suite)
	}
	if !strings.Contains(suite, "tr.Heavy") {
		t.Fatalf("suite should reference Heavy for future parallel policy:\n%s", suite)
	}
	if !strings.Contains(suite, WorkspaceAllTreesImport()) {
		t.Fatalf("suite must blank-import alltrees:\n%s", suite)
	}
}

func TestWriteWorkspaceExtrasTwoTrees(t *testing.T) {
	genRoot := t.TempDir()
	// Minimal tree packages so wreg imports resolve at compile time later.
	for _, tr := range []string{"tree-a", "tree-b"} {
		regDir := UnifiedRegistryDirForTree(genRoot, tr)
		if err := os.MkdirAll(regDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := WriteFormattedGo(filepath.Join(regDir, "registry.go"), AssembleUnifiedRegistrySource()); err != nil {
			t.Fatal(err)
		}
		allDir := UnifiedAllLeavesDirForTree(genRoot, tr)
		if err := os.MkdirAll(allDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := WriteFormattedGo(filepath.Join(allDir, "all.go"), AssembleUnifiedAllLeavesSource(nil)); err != nil {
			t.Fatal(err)
		}
		if err := WriteTreeWreg(genRoot, tr, tr == "tree-b"); err != nil {
			t.Fatal(err)
		}
	}
	if err := WriteWorkspaceExtras(genRoot, []string{"tree-a", "tree-b"}); err != nil {
		t.Fatalf("WriteWorkspaceExtras: %v", err)
	}
	for _, p := range []string{
		filepath.Join(WorkspaceRegistryDir(genRoot), "registry.go"),
		filepath.Join(WorkspaceAllTreesDir(genRoot), "all.go"),
		filepath.Join(WorkspaceSuiteDir(genRoot), "suite_test.go"),
	} {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("missing %s: %v", p, err)
		}
	}
	allBytes, err := os.ReadFile(filepath.Join(WorkspaceAllTreesDir(genRoot), "all.go"))
	if err != nil {
		t.Fatal(err)
	}
	allSrc := string(allBytes)
	if !strings.Contains(allSrc, TreeWregImportForTree("tree-a")) {
		t.Fatalf("alltrees missing tree-a wreg:\n%s", allSrc)
	}
	if !strings.Contains(allSrc, TreeWregImportForTree("tree-b")) {
		t.Fatalf("alltrees missing tree-b wreg:\n%s", allSrc)
	}

	// Compile workspace suite (empty leaf registries — suite still builds).
	modRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(modRoot, "go.mod")); err != nil {
		t.Skip("module root not found")
	}
	goMod := "module testcase\n\ngo 1.22\n\nrequire github.com/xhd2015/doctest v0.0.0\n\nreplace github.com/xhd2015/doctest => " + filepath.ToSlash(modRoot) + "\n"
	if err := os.WriteFile(filepath.Join(genRoot, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}
	// session is imported by empty registries? tree registry doesn't need session.
	// wreg imports testing/os/strings only + registries.
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = genRoot
	if out, err := tidy.CombinedOutput(); err != nil {
		t.Fatalf("tidy: %v\n%s", err, out)
	}
	cmd := exec.Command("go", "test", "./__workspace/suite/", "-count=1", "-c", "-o", os.DevNull)
	cmd.Dir = genRoot
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test -c workspace suite: %v\n%s", err, out)
	}
}
