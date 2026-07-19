package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Unified-mode package layout under gen root (module path "testcase"):
//
//	genRoot/
//	  go.mod
//	  <treeRel>/__droot/          package droot
//	  <treeRel>/__registry/       package registry — Register/All
//	  <treeRel>/__allleaves/      blank-imports every leaf package
//	  <treeRel>/<leaf>/leaf.go    non-test RunTestLeaf + init Register
//	  <treeRel>/suite/suite_test.go  iterates registry.All via t.Run
const (
	UnifiedRegistryDirName  = "__registry"
	UnifiedRegistryPkgName  = "registry"
	UnifiedAllLeavesDirName = "__allleaves"
	UnifiedAllLeavesPkgName = "allleaves"
	UnifiedSuiteDirName     = "suite"
	UnifiedSuitePkgName     = "suite"
	// RunTestLeafName is the leaf entrypoint registered with the suite.
	RunTestLeafName = "RunTestLeaf"
)

// UnifiedRegistryImportForTree returns the import path for the tree-scoped registry.
func UnifiedRegistryImportForTree(treeRel string) string {
	return treeScopedImport(treeRel, UnifiedRegistryDirName)
}

// UnifiedAllLeavesImportForTree returns the import path for __allleaves.
func UnifiedAllLeavesImportForTree(treeRel string) string {
	return treeScopedImport(treeRel, UnifiedAllLeavesDirName)
}

// UnifiedSuiteDirForTree returns the filesystem directory for suite under genRoot.
func UnifiedSuiteDirForTree(genRoot, treeRel string) string {
	return treeScopedDir(genRoot, treeRel, UnifiedSuiteDirName)
}

// UnifiedRegistryDirForTree returns the filesystem directory for __registry.
func UnifiedRegistryDirForTree(genRoot, treeRel string) string {
	return treeScopedDir(genRoot, treeRel, UnifiedRegistryDirName)
}

// UnifiedAllLeavesDirForTree returns the filesystem directory for __allleaves.
func UnifiedAllLeavesDirForTree(genRoot, treeRel string) string {
	return treeScopedDir(genRoot, treeRel, UnifiedAllLeavesDirName)
}

func treeScopedImport(treeRel, dirName string) string {
	treeRel = filepath.ToSlash(filepath.Clean(treeRel))
	if treeRel == "" || treeRel == "." {
		return "testcase/" + dirName
	}
	return "testcase/" + strings.TrimPrefix(treeRel, "./") + "/" + dirName
}

func treeScopedDir(genRoot, treeRel, dirName string) string {
	treeRel = filepath.Clean(treeRel)
	if treeRel == "" || treeRel == "." {
		return filepath.Join(genRoot, dirName)
	}
	return filepath.Join(genRoot, treeRel, dirName)
}

// LeafImportForTree returns the go import path for a leaf package under the gen module.
// leafRel is the path of the leaf gen dir relative to genRoot (filepath-style).
func LeafImportForTree(leafRel string) string {
	leafRel = filepath.ToSlash(filepath.Clean(leafRel))
	if leafRel == "" || leafRel == "." {
		return "testcase"
	}
	return "testcase/" + strings.TrimPrefix(leafRel, "./")
}

// AssembleUnifiedRegistrySource emits the tree-local registry package.
func AssembleUnifiedRegistrySource() string {
	return `package registry

import (
	"sort"
	"testing"
)

// Entry is one registered leaf test under this DOCTEST tree.
type Entry struct {
	Path   string
	Labels []string
	Fn     func(*testing.T)
}

var entries []Entry

// Register appends a leaf entry. Called from leaf package init.
func Register(e Entry) {
	entries = append(entries, e)
}

// All returns registered entries sorted by Path for stable suite order.
func All() []Entry {
	out := append([]Entry(nil), entries...)
	sort.Slice(out, func(i, j int) bool {
		return out[i].Path < out[j].Path
	})
	return out
}
`
}

// AssembleUnifiedAllLeavesSource emits a package that blank-imports every leaf.
// leafImports must be sorted for stable output.
func AssembleUnifiedAllLeavesSource(leafImports []string) string {
	var buf strings.Builder
	buf.WriteString("package ")
	buf.WriteString(UnifiedAllLeavesPkgName)
	buf.WriteString("\n\n")
	if len(leafImports) == 0 {
		return buf.String()
	}
	buf.WriteString("import (\n")
	for _, imp := range leafImports {
		buf.WriteString("\t_ \"")
		buf.WriteString(imp)
		buf.WriteString("\"\n")
	}
	buf.WriteString(")\n")
	return buf.String()
}

// AssembleUnifiedSuiteTestSource emits suite/suite_test.go that iterates All().
// Explicit aliases are required: path basenames are __registry/__allleaves while
// package names are registry/allleaves — bare imports get dropped by goimports.
func AssembleUnifiedSuiteTestSource(registryImport, allLeavesImport string) string {
	var buf strings.Builder
	buf.WriteString("package ")
	buf.WriteString(UnifiedSuitePkgName)
	buf.WriteString("\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"os\"\n")
	buf.WriteString("\t\"testing\"\n\n")
	buf.WriteString("\t")
	buf.WriteString(UnifiedRegistryPkgName)
	buf.WriteString(" \"")
	buf.WriteString(registryImport)
	buf.WriteString("\"\n")
	buf.WriteString("\t_ \"")
	buf.WriteString(allLeavesImport)
	buf.WriteString("\"\n")
	buf.WriteString(")\n\n")
	// Parent leaf via env (stdlib only) so nest timing can attribute nested RunTest
	// without importing libdoc/metrics into the generated suite module.
	buf.WriteString("func TestDoctestSuite(t *testing.T) {\n")
	buf.WriteString("\tfor _, e := range registry.All() {\n")
	buf.WriteString("\t\te := e\n")
	buf.WriteString("\t\tt.Run(e.Path, func(t *testing.T) {\n")
	buf.WriteString("\t\t\t_ = os.Setenv(\"DOCTEST_METRICS_PARENT_LEAF\", e.Path)\n")
	buf.WriteString("\t\t\tdefer os.Unsetenv(\"DOCTEST_METRICS_PARENT_LEAF\")\n")
	buf.WriteString("\t\t\te.Fn(t)\n")
	buf.WriteString("\t\t})\n")
	buf.WriteString("\t}\n")
	buf.WriteString("}\n")
	return buf.String()
}

// AssembleUnifiedLeafSource emits a non-test leaf package with RunTestLeaf + init Register.
// Body matches the thin ref leaf test body (setup/run/assert).
func AssembleUnifiedLeafSource(tc TreeCase, compileOnly bool, pkgName, docTestRoot, rootImport, rootAlias, registryImport string) (string, error) {
	if pkgName == "" {
		pkgName = "testcase"
	}
	if rootImport == "" {
		rootImport = RefRootImportPath
	}
	if rootAlias == "" {
		rootAlias = RefRootPkgName
	}
	if registryImport == "" {
		registryImport = treeScopedImport(".", UnifiedRegistryDirName)
	}

	rootDocs, leafDocs := SplitRefSetupDocs(tc.SetupFiles)
	renames := collectRootSymbolRenames(rootDocs)
	rootTypes := collectRootTypeNames(rootDocs)

	var buf strings.Builder
	buf.WriteString("package ")
	buf.WriteString(pkgName)
	buf.WriteString("\n\n")

	importsMap := collectImports(leafDocs, tc.AssertFile.GoBlock)
	for _, pkg := range []string{"testing", "syscall", sessionImportPath} {
		if _, ok := importsMap[pkg]; !ok {
			importsMap[pkg] = &ImportSpec{Path: pkg}
		}
	}
	importsMap[rootAlias+"\x00"+rootImport] = &ImportSpec{Name: rootAlias, Path: rootImport}
	importsMap[UnifiedRegistryPkgName+"\x00"+registryImport] = &ImportSpec{Name: UnifiedRegistryPkgName, Path: registryImport}
	writeImportBlock(&buf, importsMap)

	// Register this leaf with the tree suite on package init.
	buf.WriteString("func init() {\n")
	buf.WriteString("\tregistry.Register(registry.Entry{\n")
	buf.WriteString(fmt.Sprintf("\t\tPath:   %s,\n", strconv.Quote(tc.Path)))
	buf.WriteString("\t\tLabels: ")
	buf.WriteString(goStringSliceLiteral(tc.Labels))
	buf.WriteString(",\n")
	buf.WriteString("\t\tFn:     ")
	buf.WriteString(RunTestLeafName)
	buf.WriteString(",\n")
	buf.WriteString("\t})\n")
	buf.WriteString("}\n\n")

	// Leaf-only types/helpers — rewrite any references to root symbols.
	var leafBlob strings.Builder
	writePackageLevelTypesAndMethods(&leafBlob, leafDocs, tc.AssertFile.GoBlock)
	writePackageLevelConstVars(&leafBlob, leafDocs, tc.AssertFile.GoBlock)
	writePackageLevelHelpers(&leafBlob, leafDocs, tc.AssertFile.GoBlock)
	leafTop := qualifyRootSymbols(leafBlob.String(), rootAlias, renames)
	leafTop = qualifyRootTypes(leafTop, rootAlias, rootTypes)
	buf.WriteString(leafTop)

	buf.WriteString("func ")
	buf.WriteString(RunTestLeafName)
	buf.WriteString("(t *testing.T) {\n")

	writeDoctestDConstruct(&buf, docTestRoot, tc.Path)

	buf.WriteString("\tRun := ")
	buf.WriteString(rootAlias)
	buf.WriteString(".Run\n")
	buf.WriteString(fmt.Sprintf("\treq := &%s.Request{}\n", rootAlias))

	for i, doc := range rootDocs {
		if doc.GoBlock == nil || doc.GoBlock.Setup == nil {
			continue
		}
		name := fmt.Sprintf("RootSetup%d", i)
		buf.WriteString(fmt.Sprintf("\tif err := %s.%s(t, d, req); err != nil {\n", rootAlias, name))
		buf.WriteString(fmt.Sprintf("\t\tt.Fatalf(\"%s failed: %%v\", err)\n", escapeString(doc.Path)))
		buf.WriteString("\t}\n")
	}

	for i, doc := range leafDocs {
		if doc.GoBlock == nil || doc.GoBlock.Setup == nil {
			continue
		}
		name := fmt.Sprintf("setup%d", i)
		fn := *doc.GoBlock.Setup
		fn.Params = qualifyRootTypes(fn.Params, rootAlias, rootTypes)
		fn.Results = qualifyRootTypes(fn.Results, rootAlias, rootTypes)
		fn.ResultTypes = qualifyRootTypes(fn.ResultTypes, rootAlias, rootTypes)
		fn.ClosureResults = qualifyRootTypes(fn.ClosureResults, rootAlias, rootTypes)
		fn.Body = qualifyRootTypesInBody(fn.Body, rootAlias, rootTypes)
		fn.Body = qualifyRootSymbols(fn.Body, rootAlias, renames)
		fn.Params = qualifyRootSymbols(fn.Params, rootAlias, renames)
		fn.Params = ensureDoctestParam(fn.Params)
		writeFuncClosure(&buf, name, fn)
		buf.WriteString(fmt.Sprintf("\tif err := %s(t, d, req); err != nil {\n", name))
		buf.WriteString(fmt.Sprintf("\t\tt.Fatalf(\"%s failed: %%v\", err)\n", escapeString(doc.Path)))
		buf.WriteString("\t}\n")
	}

	assertFn := *tc.AssertFile.GoBlock.Assert
	assertFn.Params = qualifyRootTypes(assertFn.Params, rootAlias, rootTypes)
	assertFn.Results = qualifyRootTypes(assertFn.Results, rootAlias, rootTypes)
	assertFn.ResultTypes = qualifyRootTypes(assertFn.ResultTypes, rootAlias, rootTypes)
	assertFn.ClosureResults = qualifyRootTypes(assertFn.ClosureResults, rootAlias, rootTypes)
	assertFn.Body = qualifyRootTypesInBody(assertFn.Body, rootAlias, rootTypes)
	assertFn.Body = qualifyRootSymbols(assertFn.Body, rootAlias, renames)
	assertFn.Params = ensureDoctestParam(assertFn.Params)
	writeFuncClosure(&buf, "assert", assertFn)

	buf.WriteString("\t_ = Run\n")
	helperNames := collectHelperNames(leafDocs, tc.AssertFile.GoBlock)
	for _, name := range helperNames {
		buf.WriteString(fmt.Sprintf("\t_ = %s\n", name))
	}

	if compileOnly {
		buf.WriteString("\t// compileOnly\n")
		buf.WriteString("\t_ = d\n")
		buf.WriteString("\t_ = req\n")
		buf.WriteString("\t_ = assert\n")
		buf.WriteString(fmt.Sprintf("\tvar resp *%s.Response\n", rootAlias))
		buf.WriteString("\tvar runErr error\n")
		buf.WriteString("\t_ = resp\n")
		buf.WriteString("\t_ = runErr\n")
		buf.WriteString("}\n")
		return buf.String(), nil
	}

	buf.WriteString(fmt.Sprintf("\tresp, runErr := %s.Run(t, d, req)\n", rootAlias))
	buf.WriteString("\tassert(t, d, req, resp, runErr)\n")
	buf.WriteString("}\n")
	return buf.String(), nil
}

func goStringSliceLiteral(ss []string) string {
	if len(ss) == 0 {
		return "nil"
	}
	var b strings.Builder
	b.WriteString("[]string{")
	for i, s := range ss {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(strconv.Quote(s))
	}
	b.WriteString("}")
	return b.String()
}

// UnifiedLeafFileName returns the non-test .go basename for a unified leaf package.
// Always "leaf.go" — never "*_test.go", because Go treats a trailing _test.go
// suffix as a test file (e.g. path "suite-only-go-test" → suite_only_go_test.go
// would leave the package with "no non-test Go files").
func UnifiedLeafFileName(tc TreeCase) string {
	_ = tc
	return "leaf.go"
}

// WriteUnifiedLeafCase writes a non-test leaf package under leafDir.
func WriteUnifiedLeafCase(leafDir string, tc TreeCase, compileOnly bool, pkgName, docTestRoot, rootImport, registryImport string) (string, error) {
	if err := os.MkdirAll(leafDir, 0755); err != nil {
		return "", err
	}
	if rootImport == "" {
		rootImport = RefRootImportPath
	}
	src, err := AssembleUnifiedLeafSource(tc, compileOnly, pkgName, docTestRoot, rootImport, RefRootPkgName, registryImport)
	if err != nil {
		return "", fmt.Errorf("%s: %w", tc.Path, err)
	}
	leafPath := filepath.Join(leafDir, UnifiedLeafFileName(tc))
	if err := WriteFormattedGo(leafPath, src); err != nil {
		return "", err
	}
	return leafPath, nil
}

// WriteUnifiedTreeExtras writes __registry, __allleaves, and suite for a tree.
// leafImports are full module import paths for each leaf package (already sorted).
func WriteUnifiedTreeExtras(genRoot, treeRel string, leafImports []string) error {
	regDir := UnifiedRegistryDirForTree(genRoot, treeRel)
	if err := os.MkdirAll(regDir, 0755); err != nil {
		return err
	}
	if err := WriteFormattedGo(filepath.Join(regDir, "registry.go"), AssembleUnifiedRegistrySource()); err != nil {
		return fmt.Errorf("write unified registry: %w", err)
	}

	allDir := UnifiedAllLeavesDirForTree(genRoot, treeRel)
	if err := os.MkdirAll(allDir, 0755); err != nil {
		return err
	}
	sorted := append([]string(nil), leafImports...)
	sort.Strings(sorted)
	if err := WriteFormattedGo(filepath.Join(allDir, "all.go"), AssembleUnifiedAllLeavesSource(sorted)); err != nil {
		return fmt.Errorf("write unified allleaves: %w", err)
	}

	suiteDir := UnifiedSuiteDirForTree(genRoot, treeRel)
	if err := os.MkdirAll(suiteDir, 0755); err != nil {
		return err
	}
	suiteSrc := AssembleUnifiedSuiteTestSource(
		UnifiedRegistryImportForTree(treeRel),
		UnifiedAllLeavesImportForTree(treeRel),
	)
	if err := WriteFormattedGo(filepath.Join(suiteDir, "suite_test.go"), suiteSrc); err != nil {
		return fmt.Errorf("write unified suite: %w", err)
	}
	return nil
}
