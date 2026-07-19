package core

import (
	"crypto/sha256"
	"encoding/hex"
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
//	  <treeRel>/__droot/             package droot
//	  <treeRel>/<parent>/setup.go    intermediate packages (shared, same as ref)
//	  <treeRel>/__registry/          package registry — Register/All
//	  <treeRel>/__allleaves/         blank-imports every leaf package
//	  <treeRel>/<leaf>/leaf.go       non-test RunTestLeaf + init Register
//	  <treeRel>/suite/suite_test.go  iterates registry.All via t.Run
//
// Intermediate SETUP packages are written once via WriteRefIntermediatePackages
// (shared with ref mode). Leaves import ancestors and call RootSetup* → Setup → leaf setups.
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
//
// leavesFingerprint is a content hash of all leaf package sources. It is emitted
// as a package-level const so any leaf SETUP/ASSERT rewrite changes suite source
// and thus the suite test binary identity. Blank-import alone is not enough for
// go test result cache: leaf packages reached only via init/function pointers
// can leave the binary content-id stable from the cache key’s point of view.
func AssembleUnifiedSuiteTestSource(registryImport, allLeavesImport, leavesFingerprint string) string {
	if leavesFingerprint == "" {
		leavesFingerprint = "empty"
	}
	var buf strings.Builder
	buf.WriteString("package ")
	buf.WriteString(UnifiedSuitePkgName)
	buf.WriteString("\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"os\"\n")
	buf.WriteString("\t\"strings\"\n")
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
	// Force suite package rebuild whenever any leaf body changes.
	// Use a package-level var (not const): go may fold/eliminate unused consts
	// such that the suite test binary content-id stays stable and go test
	// falsely reports (cached) after leaf rewrites.
	buf.WriteString("// Code generated by doctest; DO NOT EDIT.\n")
	buf.WriteString("// fingerprint of all leaf package sources (cache invalidation).\n")
	buf.WriteString("var doctestLeavesFingerprint = \"")
	buf.WriteString(leavesFingerprint)
	buf.WriteString("\"\n\n")
	// Parent leaf via env (stdlib only) so nest timing can attribute nested RunTest
	// without importing libdoc/metrics into the generated suite module.
	//
	// Subtest names use "__" instead of "/" so go test -json does not invent
	// intermediate path nodes (t.Run("a/b") creates TestDoctestSuite/a and
	// …/a/b, which double-counts and breaks multi-leaf summaries).
	buf.WriteString("func TestDoctestSuite(t *testing.T) {\n")
	buf.WriteString("\tif len(doctestLeavesFingerprint) == 0 {\n")
	buf.WriteString("\t\tt.Fatal(\"empty doctestLeavesFingerprint\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tfor _, e := range registry.All() {\n")
	buf.WriteString("\t\te := e\n")
	buf.WriteString("\t\t// Encode \"/\" as \"__\" so go test does not nest path segments.\n")
	buf.WriteString("\t\tname := strings.ReplaceAll(e.Path, \"/\", \"__\")\n")
	buf.WriteString("\t\tt.Run(name, func(t *testing.T) {\n")
	buf.WriteString("\t\t\t_ = os.Setenv(\"DOCTEST_METRICS_PARENT_LEAF\", e.Path)\n")
	buf.WriteString("\t\t\tdefer os.Unsetenv(\"DOCTEST_METRICS_PARENT_LEAF\")\n")
	buf.WriteString("\t\t\te.Fn(t)\n")
	buf.WriteString("\t\t})\n")
	buf.WriteString("\t}\n")
	buf.WriteString("}\n")
	return buf.String()
}

// AssembleUnifiedLeafSource emits a non-test leaf package with RunTestLeaf + init Register.
// Body mirrors hierarchical AssembleRefLeafTestSource: leaf-local docs only, import droot
// and intermediate packages, call RootSetup* → intermediate Setup → leaf setups → Run → assert.
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

	part := PartitionRefSetupDocs(tc)
	rootDocs := part.RootDocs
	leafDocs := part.LeafDocs


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
	for _, g := range part.Intermediate {
		alias := RefIntermediateAlias(g.Dir)
		imp := RefIntermediateImport(rootImport, g.Dir)
		importsMap[alias+"\x00"+imp] = &ImportSpec{Name: alias, Path: imp}
	}
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

	// Leaf-only types/helpers — rewrite references to root + intermediate symbols.
	var leafBlob strings.Builder
	writePackageLevelTypesAndMethods(&leafBlob, leafDocs, tc.AssertFile.GoBlock)
	writePackageLevelConstVars(&leafBlob, leafDocs, tc.AssertFile.GoBlock)
	writePackageLevelHelpers(&leafBlob, leafDocs, tc.AssertFile.GoBlock)
	leafTop := qualifyAncestorSymbols(leafBlob.String(), part, rootAlias, rootDocs)
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

	// Intermediate setups (exported package funcs), parents first.
	for _, g := range part.Intermediate {
		alias := RefIntermediateAlias(g.Dir)
		setupTotal := 0
		for _, doc := range g.Docs {
			if doc.GoBlock != nil && doc.GoBlock.Setup != nil {
				setupTotal++
			}
		}
		setupIdx := 0
		for _, doc := range g.Docs {
			if doc.GoBlock == nil || doc.GoBlock.Setup == nil {
				continue
			}
			name := intermediateSetupName(setupIdx, setupTotal)
			setupIdx++
			buf.WriteString(fmt.Sprintf("\tif err := %s.%s(t, d, req); err != nil {\n", alias, name))
			buf.WriteString(fmt.Sprintf("\t\tt.Fatalf(\"%s failed: %%v\", err)\n", escapeString(doc.Path)))
			buf.WriteString("\t}\n")
		}
	}

	for i, doc := range leafDocs {
		if doc.GoBlock == nil || doc.GoBlock.Setup == nil {
			continue
		}
		name := fmt.Sprintf("setup%d", i)
		fn := *doc.GoBlock.Setup
		fn.Params = qualifyAncestorSymbols(fn.Params, part, rootAlias, rootDocs)
		fn.Results = qualifyAncestorSymbols(fn.Results, part, rootAlias, rootDocs)
		fn.ResultTypes = qualifyAncestorSymbols(fn.ResultTypes, part, rootAlias, rootDocs)
		fn.ClosureResults = qualifyAncestorSymbols(fn.ClosureResults, part, rootAlias, rootDocs)
		fn.Body = qualifyAncestorSymbols(fn.Body, part, rootAlias, rootDocs)
		fn.Params = ensureDoctestParam(fn.Params)
		writeFuncClosure(&buf, name, fn)
		buf.WriteString(fmt.Sprintf("\tif err := %s(t, d, req); err != nil {\n", name))
		buf.WriteString(fmt.Sprintf("\t\tt.Fatalf(\"%s failed: %%v\", err)\n", escapeString(doc.Path)))
		buf.WriteString("\t}\n")
	}

	assertFn := *tc.AssertFile.GoBlock.Assert
	assertFn.Params = qualifyAncestorSymbols(assertFn.Params, part, rootAlias, rootDocs)
	assertFn.Results = qualifyAncestorSymbols(assertFn.Results, part, rootAlias, rootDocs)
	assertFn.ResultTypes = qualifyAncestorSymbols(assertFn.ResultTypes, part, rootAlias, rootDocs)
	assertFn.ClosureResults = qualifyAncestorSymbols(assertFn.ClosureResults, part, rootAlias, rootDocs)
	assertFn.Body = qualifyAncestorSymbols(assertFn.Body, part, rootAlias, rootDocs)
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
// Returns the written path and the unformatted source (for suite fingerprinting).
func WriteUnifiedLeafCase(leafDir string, tc TreeCase, compileOnly bool, pkgName, docTestRoot, rootImport, registryImport string) (path, src string, err error) {
	if err := os.MkdirAll(leafDir, 0755); err != nil {
		return "", "", err
	}
	if rootImport == "" {
		rootImport = RefRootImportPath
	}
	src, err = AssembleUnifiedLeafSource(tc, compileOnly, pkgName, docTestRoot, rootImport, RefRootPkgName, registryImport)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", tc.Path, err)
	}
	leafPath := filepath.Join(leafDir, UnifiedLeafFileName(tc))
	if err := WriteFormattedGo(leafPath, src); err != nil {
		return "", "", err
	}
	return leafPath, src, nil
}

// WriteUnifiedTree writes __droot, intermediate packages once, unified leaves,
// __registry, __allleaves, and suite under genRoot (flat treeRel ".").
// Useful for unit tests; production generation uses generateContext.writeUnifiedCases.
func WriteUnifiedTree(genRoot string, cases []TreeCase, docTestRoot string, compileOnly bool, pkgName string) error {
	if len(cases) == 0 {
		return fmt.Errorf("WriteUnifiedTree: no cases")
	}
	if pkgName == "" {
		pkgName = "testcase"
	}

	rootDocs, _ := SplitRefSetupDocs(cases[0].SetupFiles)
	if len(rootDocs) == 0 {
		rootDocs = cases[0].SetupFiles
	}
	rootSrc, err := AssembleRefRootSource(rootDocs, RefRootPkgName)
	if err != nil {
		return err
	}
	rootDir := filepath.Join(genRoot, RefRootDirName)
	rootImport := RefRootImportPath
	registryImport := UnifiedRegistryImportForTree(".")
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		return err
	}
	if err := WriteFormattedGo(filepath.Join(rootDir, "droot.go"), rootSrc); err != nil {
		return fmt.Errorf("write unified ref root package: %w", err)
	}

	if err := WriteRefIntermediatePackages(genRoot, ".", rootImport, rootDocs, cases); err != nil {
		return err
	}

	leafImports := make([]string, 0, len(cases))
	leafSources := make(map[string]string, len(cases))
	for _, tc := range cases {
		leafDir := genRoot
		if tc.Path != "" {
			leafDir = filepath.Join(genRoot, tc.Path)
		}
		_, src, err := WriteUnifiedLeafCase(leafDir, tc, compileOnly, pkgName, docTestRoot, rootImport, registryImport)
		if err != nil {
			return err
		}
		leafRel := tc.Path
		if leafRel == "" {
			leafRel = "."
		}
		imp := LeafImportForTree(leafRel)
		leafImports = append(leafImports, imp)
		leafSources[imp] = src
	}

	return WriteUnifiedTreeExtras(genRoot, ".", leafImports, leafSources)
}

// WriteUnifiedTreeExtras writes __registry, __allleaves, and suite for a tree.
// leafImports are full module import paths for each leaf package (already sorted).
// leafSources maps import path → leaf.go source (used for suite fingerprint).
func WriteUnifiedTreeExtras(genRoot, treeRel string, leafImports []string, leafSources map[string]string) error {
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
	fp := hashLeafSources(sorted, leafSources)
	suiteSrc := AssembleUnifiedSuiteTestSource(
		UnifiedRegistryImportForTree(treeRel),
		UnifiedAllLeavesImportForTree(treeRel),
		fp,
	)
	if err := WriteFormattedGo(filepath.Join(suiteDir, "suite_test.go"), suiteSrc); err != nil {
		return fmt.Errorf("write unified suite: %w", err)
	}
	return nil
}

// hashLeafSources returns a stable hex fingerprint of leaf package sources.
func hashLeafSources(sortedImports []string, leafSources map[string]string) string {
	h := sha256.New()
	for _, imp := range sortedImports {
		h.Write([]byte(imp))
		h.Write([]byte{0})
		if src, ok := leafSources[imp]; ok {
			h.Write([]byte(src))
		}
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))[:16]
}
