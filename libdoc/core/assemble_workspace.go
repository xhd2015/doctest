package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Workspace layout under the shared gen module (module path "testcase"):
//
//	genRoot/
//	  __workspace/
//	    __registry/   package registry — TreeEntry Register/All
//	    __alltrees/   blank-imports each selected tree's __wreg
//	    suite/        TestDoctestSuite iterates workspace registry
//	  <treeRel>/__wreg/  init registers this tree into workspace registry
//
// Leaf → tree registry is unchanged. Tree → workspace mirrors that pattern.

const (
	WorkspaceDirName           = "__workspace"
	WorkspaceRegistryDirName   = "__registry"
	WorkspaceRegistryPkgName   = "registry"
	WorkspaceAllTreesDirName   = "__alltrees"
	WorkspaceAllTreesPkgName   = "alltrees"
	WorkspaceSuiteDirName      = "suite"
	WorkspaceSuitePkgName      = "suite"
	// TreeWregDirName is the per-tree package that Registers into the workspace.
	TreeWregDirName = "__wreg"
	TreeWregPkgName = "wreg"
)

// Workspace paths (import + filesystem).

func WorkspaceRootDir(genRoot string) string {
	return filepath.Join(genRoot, WorkspaceDirName)
}

func WorkspaceRegistryDir(genRoot string) string {
	return filepath.Join(genRoot, WorkspaceDirName, WorkspaceRegistryDirName)
}

func WorkspaceAllTreesDir(genRoot string) string {
	return filepath.Join(genRoot, WorkspaceDirName, WorkspaceAllTreesDirName)
}

func WorkspaceSuiteDir(genRoot string) string {
	return filepath.Join(genRoot, WorkspaceDirName, WorkspaceSuiteDirName)
}

func WorkspaceRegistryImport() string {
	return "testcase/" + WorkspaceDirName + "/" + WorkspaceRegistryDirName
}

func WorkspaceAllTreesImport() string {
	return "testcase/" + WorkspaceDirName + "/" + WorkspaceAllTreesDirName
}

func WorkspaceSuiteImport() string {
	return "testcase/" + WorkspaceDirName + "/" + WorkspaceSuiteDirName
}

func TreeWregDirForTree(genRoot, treeRel string) string {
	return treeScopedDir(genRoot, treeRel, TreeWregDirName)
}

func TreeWregImportForTree(treeRel string) string {
	return treeScopedImport(treeRel, TreeWregDirName)
}

// AssembleWorkspaceRegistrySource emits the workspace-level tree registry.
func AssembleWorkspaceRegistrySource() string {
	return `package registry

import (
	"sort"
	"testing"
)

// TreeEntry is one DOCTEST root registered under the workspace suite.
type TreeEntry struct {
	Path string // treeRel relative to module root (e.g. "libdoc/build/tests")
	// Run executes all leaves for this tree (typically via the tree registry).
	Run func(*testing.T)
}

var entries []TreeEntry

// Register appends a tree entry. Called from each tree's __wreg init.
func Register(e TreeEntry) {
	entries = append(entries, e)
}

// All returns registered trees sorted by Path for stable suite order.
func All() []TreeEntry {
	out := append([]TreeEntry(nil), entries...)
	sort.Slice(out, func(i, j int) bool {
		return out[i].Path < out[j].Path
	})
	return out
}
`
}

// AssembleTreeWregSource emits <treeRel>/__wreg that registers into the workspace.
// treeRelSlash is the slash-separated tree path stored on TreeEntry.Path (may be ".").
// Every leaf subtest calls t.Parallel().
func AssembleTreeWregSource(treeRelSlash string, treeRegistryImport, treeAllLeavesImport, workspaceRegistryImport string) string {
	if treeRelSlash == "" {
		treeRelSlash = "."
	}
	var buf strings.Builder
	buf.WriteString("package ")
	buf.WriteString(TreeWregPkgName)
	buf.WriteString("\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"strings\"\n")
	buf.WriteString("\t\"syscall\"\n")
	buf.WriteString("\t\"testing\"\n\n")
	buf.WriteString("\ttreereg \"")
	buf.WriteString(treeRegistryImport)
	buf.WriteString("\"\n")
	buf.WriteString("\t_ \"")
	buf.WriteString(treeAllLeavesImport)
	buf.WriteString("\"\n")
	buf.WriteString("\twsreg \"")
	buf.WriteString(workspaceRegistryImport)
	buf.WriteString("\"\n")
	buf.WriteString(")\n\n")
	buf.WriteString("func init() {\n")
	buf.WriteString("\twsreg.Register(wsreg.TreeEntry{\n")
	buf.WriteString("\t\tPath: ")
	buf.WriteString(strconvQuote(treeRelSlash))
	buf.WriteString(",\n")
	// Leaf nest parent is d.Metrics.ParentLeaf inside RunTestLeaf — no process env.
	// Warm leaf-cache skips use DOCTEST_LEAF_CACHE_SKIP_PATHS (same as tree suite RunAll).
	buf.WriteString("\t\tRun: func(t *testing.T) {\n")
	buf.WriteString("\t\t\tskip := map[string]struct{}{}\n")
	buf.WriteString("\t\t\tif s, ok := syscall.Getenv(\"DOCTEST_LEAF_CACHE_SKIP_PATHS\"); ok && s != \"\" {\n")
	buf.WriteString("\t\t\t\tfor _, p := range strings.Split(s, \"\\n\") {\n")
	buf.WriteString("\t\t\t\t\tp = strings.TrimSpace(p)\n")
	buf.WriteString("\t\t\t\t\tif p != \"\" {\n")
	buf.WriteString("\t\t\t\t\t\tskip[p] = struct{}{}\n")
	buf.WriteString("\t\t\t\t\t}\n")
	buf.WriteString("\t\t\t\t}\n")
	buf.WriteString("\t\t\t}\n")
	buf.WriteString("\t\t\tfor _, e := range treereg.All() {\n")
	buf.WriteString("\t\t\t\te := e\n")
	buf.WriteString("\t\t\t\t// Encode \"/\" as \"__\" so go test does not nest path segments.\n")
	buf.WriteString("\t\t\t\tname := strings.ReplaceAll(e.Path, \"/\", \"__\")\n")
	buf.WriteString("\t\t\t\tt.Run(name, func(t *testing.T) {\n")
	buf.WriteString("\t\t\t\t\tt.Parallel()\n")
	buf.WriteString("\t\t\t\t\tif _, hit := skip[e.Path]; hit {\n")
	buf.WriteString("\t\t\t\t\t\treturn\n")
	buf.WriteString("\t\t\t\t\t}\n")
	buf.WriteString("\t\t\t\t\te.Fn(t)\n")
	buf.WriteString("\t\t\t\t})\n")
	buf.WriteString("\t\t\t}\n")
	buf.WriteString("\t\t},\n")
	buf.WriteString("\t})\n")
	buf.WriteString("}\n")
	return buf.String()
}

func strconvQuote(s string) string {
	// Minimal quoted string for generated Go source.
	return fmt.Sprintf("%q", s)
}

// AssembleWorkspaceAllTreesSource blank-imports every selected tree __wreg.
// wregImports must be sorted for stable output.
func AssembleWorkspaceAllTreesSource(wregImports []string) string {
	var buf strings.Builder
	buf.WriteString("package ")
	buf.WriteString(WorkspaceAllTreesPkgName)
	buf.WriteString("\n\n")
	if len(wregImports) == 0 {
		return buf.String()
	}
	buf.WriteString("import (\n")
	for _, imp := range wregImports {
		buf.WriteString("\t_ \"")
		buf.WriteString(imp)
		buf.WriteString("\"\n")
	}
	buf.WriteString(")\n")
	return buf.String()
}

// AssembleWorkspaceRunAllSource emits __workspace/suite/runall.go (importable).
// Every tree subtest calls t.Parallel(); leaves Parallel inside each tree Run.
func AssembleWorkspaceRunAllSource(workspaceRegistryImport, allTreesImport string) string {
	var buf strings.Builder
	buf.WriteString("package ")
	buf.WriteString(WorkspaceSuitePkgName)
	buf.WriteString("\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"strings\"\n")
	buf.WriteString("\t\"testing\"\n\n")
	buf.WriteString("\t")
	buf.WriteString(WorkspaceRegistryPkgName)
	buf.WriteString(" \"")
	buf.WriteString(workspaceRegistryImport)
	buf.WriteString("\"\n")
	buf.WriteString("\t_ \"")
	buf.WriteString(allTreesImport)
	buf.WriteString("\"\n")
	buf.WriteString(")\n\n")
	buf.WriteString("// RunAll runs every registered tree (importable from a multi-module hub).\n")
	buf.WriteString("func RunAll(t *testing.T) {\n")
	buf.WriteString("\tfor _, tr := range registry.All() {\n")
	buf.WriteString("\t\ttr := tr\n")
	buf.WriteString("\t\t// Encode \"/\" as \"__\" so go test does not invent intermediate nodes.\n")
	buf.WriteString("\t\tname := strings.ReplaceAll(tr.Path, \"/\", \"__\")\n")
	buf.WriteString("\t\tt.Run(name, func(t *testing.T) {\n")
	buf.WriteString("\t\t\tt.Parallel()\n")
	buf.WriteString("\t\t\ttr.Run(t)\n")
	buf.WriteString("\t\t})\n")
	buf.WriteString("\t}\n")
	buf.WriteString("}\n")
	return buf.String()
}

// AssembleWorkspaceSuiteTestSource emits __workspace/suite/suite_test.go → RunAll.
func AssembleWorkspaceSuiteTestSource(workspaceRegistryImport, allTreesImport string) string {
	_ = workspaceRegistryImport
	_ = allTreesImport
	var buf strings.Builder
	buf.WriteString("package ")
	buf.WriteString(WorkspaceSuitePkgName)
	buf.WriteString("\n\n")
	buf.WriteString("import \"testing\"\n\n")
	buf.WriteString("func TestDoctestSuite(t *testing.T) {\n")
	buf.WriteString("\tRunAll(t)\n")
	buf.WriteString("}\n")
	return buf.String()
}

// WriteTreeWreg writes <treeRel>/__wreg for workspace registration.
func WriteTreeWreg(genRoot, treeRel string) error {
	treeRelSlash := filepath.ToSlash(filepath.Clean(treeRel))
	if treeRelSlash == "" {
		treeRelSlash = "."
	}
	dir := TreeWregDirForTree(genRoot, treeRel)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	src := AssembleTreeWregSource(
		treeRelSlash,
		UnifiedRegistryImportForTree(treeRel),
		UnifiedAllLeavesImportForTree(treeRel),
		WorkspaceRegistryImport(),
	)
	return WriteFormattedGo(filepath.Join(dir, "wreg.go"), src)
}

// WriteWorkspaceExtras writes workspace __registry, __alltrees, and suite for the
// selected treeRels only (fan-in list). treeRels use filesystem-relative form.
func WriteWorkspaceExtras(genRoot string, treeRels []string) error {
	if err := os.MkdirAll(WorkspaceRegistryDir(genRoot), 0755); err != nil {
		return err
	}
	if err := WriteFormattedGo(
		filepath.Join(WorkspaceRegistryDir(genRoot), "registry.go"),
		AssembleWorkspaceRegistrySource(),
	); err != nil {
		return fmt.Errorf("write workspace registry: %w", err)
	}

	wregImports := make([]string, 0, len(treeRels))
	for _, tr := range treeRels {
		wregImports = append(wregImports, TreeWregImportForTree(tr))
	}
	sort.Strings(wregImports)

	if err := os.MkdirAll(WorkspaceAllTreesDir(genRoot), 0755); err != nil {
		return err
	}
	if err := WriteFormattedGo(
		filepath.Join(WorkspaceAllTreesDir(genRoot), "all.go"),
		AssembleWorkspaceAllTreesSource(wregImports),
	); err != nil {
		return fmt.Errorf("write workspace alltrees: %w", err)
	}

	if err := os.MkdirAll(WorkspaceSuiteDir(genRoot), 0755); err != nil {
		return err
	}
	regImp := WorkspaceRegistryImport()
	allImp := WorkspaceAllTreesImport()
	suiteDir := WorkspaceSuiteDir(genRoot)
	if err := WriteFormattedGo(filepath.Join(suiteDir, "runall.go"), AssembleWorkspaceRunAllSource(regImp, allImp)); err != nil {
		return fmt.Errorf("write workspace RunAll: %w", err)
	}
	if err := WriteFormattedGo(filepath.Join(suiteDir, "suite_test.go"), AssembleWorkspaceSuiteTestSource(regImp, allImp)); err != nil {
		return fmt.Errorf("write workspace suite: %w", err)
	}
	return nil
}
