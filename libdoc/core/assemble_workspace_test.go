package core

import (
	"os"
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
	if strings.Contains(reg, "Heavy") {
		t.Fatalf("registry must not include Heavy field:\n%s", reg)
	}
	if strings.Contains(reg, "IsolateProcessEnv") {
		t.Fatalf("registry must not include IsolateProcessEnv (process env mutation):\n%s", reg)
	}

	wreg := AssembleTreeWregSource(
		"libdoc/build/tests",
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
	if !strings.Contains(wreg, "t.Parallel()") {
		t.Fatalf("wreg must Parallel leaves:\n%s", wreg)
	}
	if strings.Contains(wreg, "IsolateProcessEnv") || strings.Contains(wreg, "Setenv") {
		t.Fatalf("wreg must not mutate process env:\n%s", wreg)
	}
	if !strings.Contains(wreg, `Path: "libdoc/build/tests"`) && !strings.Contains(wreg, `Path:  "libdoc/build/tests"`) {
		// generated may use Path: "..." with single space after colon
		if !strings.Contains(wreg, `"libdoc/build/tests"`) {
			t.Fatalf("wreg path:\n%s", wreg)
		}
	}

	runAll := AssembleWorkspaceRunAllSource(WorkspaceRegistryImport(), WorkspaceAllTreesImport())
	if !strings.Contains(runAll, "registry.All()") {
		t.Fatalf("RunAll must use workspace All:\n%s", runAll)
	}
	if !strings.Contains(runAll, "t.Parallel()") {
		t.Fatalf("RunAll must Parallel trees:\n%s", runAll)
	}
	if strings.Contains(runAll, "Heavy") || strings.Contains(runAll, "tr.Heavy") {
		t.Fatalf("RunAll must not reference Heavy:\n%s", runAll)
	}
	if !strings.Contains(runAll, WorkspaceAllTreesImport()) {
		t.Fatalf("RunAll must blank-import alltrees:\n%s", runAll)
	}
	suite := AssembleWorkspaceSuiteTestSource(WorkspaceRegistryImport(), WorkspaceAllTreesImport())
	if !strings.Contains(suite, "RunAll(t)") {
		t.Fatalf("suite_test must call RunAll:\n%s", suite)
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
		if err := WriteTreeWreg(genRoot, tr); err != nil {
			t.Fatal(err)
		}
	}
	if err := WriteWorkspaceExtras(genRoot, []string{"tree-a", "tree-b"}); err != nil {
		t.Fatal(err)
	}
	// Workspace suite must compile enough to exist.
	suite := filepath.Join(WorkspaceSuiteDir(genRoot), "runall.go")
	data, err := os.ReadFile(suite)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "t.Parallel()") {
		t.Fatalf("workspace runall missing Parallel:\n%s", data)
	}
}
