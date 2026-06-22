package core

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/xhd2015/doctest/libdoc/rules"
	"golang.org/x/tools/imports"
)

func DiscoverTreeCases(root string) ([]TreeCase, error) {
	return discoverTreeCasesInternal(root, nil)
}

func DiscoverTreeCasesVerbose(root string, w io.Writer) ([]TreeCase, error) {
	return discoverTreeCasesInternal(root, w)
}

func discoverTreeCasesInternal(root string, w io.Writer) ([]TreeCase, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", root)
	}

	var verrs []ValidationError

	doctestPath := filepath.Join(root, "DOCTEST.md")
	doctestContent, err := os.ReadFile(doctestPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	doctestDoc, err := ParseDOCTESTDocument(doctestPath, string(doctestContent))
	if err != nil {
		verrs = append(verrs, ValidationError{Path: "DOCTEST.md", Msg: err.Error()})
	} else if doctestDoc.GoBlock == nil {
		verrs = append(verrs, ValidationError{Path: "DOCTEST.md", Msg: "must have a Go code block"})
	} else {
		if v := rules.CheckRootHasRequestResponse(doctestDoc.GoBlock.Types, "DOCTEST.md"); v != nil {
			verrs = append(verrs, ValidationError{Path: v.Path, Msg: v.Msg})
		}
		if v := rules.CheckRootHasRun(doctestDoc.GoBlock.Run != nil, "DOCTEST.md"); v != nil {
			verrs = append(verrs, ValidationError{Path: v.Path, Msg: v.Msg})
		}
	}

	rootSetupPath := filepath.Join(root, "SETUP.md")
	rootSetupContent, rootSetupErr := os.ReadFile(rootSetupPath)
	if rootSetupErr == nil {
		rootSetup, parseErr := ParseSetupDocument(rootSetupPath, string(rootSetupContent))
		if parseErr != nil {
			verrs = append(verrs, ValidationError{Path: "SETUP.md", Msg: parseErr.Error()})
		} else if rootSetup.GoBlock != nil {
			if v := rules.CheckRootSetupNoRequestResponseRun(rootSetup.GoBlock.Types, rootSetup.GoBlock.Run != nil, "SETUP.md"); v != nil {
				verrs = append(verrs, ValidationError{Path: v.Path, Msg: v.Msg})
			}
		}
		if w != nil {
			printSetupVerbose(w, rootSetup, "SETUP.md")
		}
	} else if !os.IsNotExist(rootSetupErr) {
		return nil, rootSetupErr
	}

	if w != nil {
		printSetupVerbose(w, doctestDoc, "DOCTEST.md")
	}

	printedSetupDirs := make(map[string]bool)
	printedSetupDirs[root] = true
	var cases []TreeCase
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			verrs = append(verrs, ValidationError{Path: path, Msg: walkErr.Error()})
			return nil
		}
		if d.IsDir() {
			if d.Name() == "testdata" {
				return filepath.SkipDir
			}
			if path == root {
				return nil
			}
			if _, err := os.Stat(filepath.Join(path, "DOCTEST.md")); err == nil {
				return filepath.SkipDir
			}
			relPath, _ := filepath.Rel(root, path)
			if w != nil {
				setupPath := filepath.Join(path, "SETUP.md")
				if !printedSetupDirs[setupPath] {
					doc, readErr := readSetup(setupPath)
					if readErr == nil && doc.GoBlock != nil {
						printedSetupDirs[setupPath] = true
						printSetupVerbose(w, doc, filepath.Join(relPath, "SETUP.md"))
					}
				}
			}
			setupPath := filepath.Join(path, "SETUP.md")
			doc, readErr := readSetup(setupPath)
			if readErr != nil {
				rel, _ := filepath.Rel(root, setupPath)
				verrs = append(verrs, ValidationError{Path: rel, Msg: readErr.Error()})
			} else if doc.GoBlock == nil {
				rel, _ := filepath.Rel(root, setupPath)
				verrs = append(verrs, ValidationError{Path: rel, Msg: "must have a Go code block"})
		} else if doc.GoBlock.Setup == nil {
			rel, _ := filepath.Rel(root, setupPath)
			verrs = append(verrs, ValidationError{Path: rel, Msg: "must have func Setup"})
			}
			return nil
		}
		if d.Name() != "ASSERT.md" {
			return nil
		}
		leafDir := filepath.Dir(path)
		relLeaf, err := filepath.Rel(root, leafDir)
		if err != nil {
			verrs = append(verrs, ValidationError{Path: path, Msg: err.Error()})
			return nil
		}
		if relLeaf == "." {
			relLeaf = ""
		}
		setupDocs, chainErr := setupChain(root, leafDir, doctestDoc)
		if chainErr != nil {
			relAssert, _ := filepath.Rel(root, path)
			verrs = append(verrs, ValidationError{Path: relAssert, Msg: chainErr.Error()})
			return nil
		}
		assertContent, err := os.ReadFile(path)
		if err != nil {
			relAssert, _ := filepath.Rel(root, path)
			verrs = append(verrs, ValidationError{Path: relAssert, Msg: err.Error()})
			return nil
		}
		relAssert, _ := filepath.Rel(root, path)
		assertDoc, err := ParseAssertDocument(relAssert, string(assertContent))
		if err != nil {
			verrs = append(verrs, ValidationError{Path: relAssert, Msg: err.Error()})
			return nil
		}
		if v := rules.CheckChainHasRun(runSource(setupDocs), relAssert); v != nil {
			verrs = append(verrs, ValidationError{Path: v.Path, Msg: v.Msg})
		}
		tc := TreeCase{
			Name:       CaseName(relLeaf),
			Path:       relLeaf,
			SetupFiles: setupDocs,
			AssertFile: assertDoc,
		}
		cases = append(cases, tc)

		if w != nil {
			fmt.Fprintf(w, "  ✦ %-30s (Run: %s)\n", CaseName(relLeaf), runSource(setupDocs))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(verrs) > 0 {
		return nil, JoinValidationErrors(verrs)
	}
	sort.Slice(cases, func(i, j int) bool { return cases[i].Path < cases[j].Path })
	return cases, nil
}

func JoinValidationErrors(verrs []ValidationError) error {
	var msgs []string
	for _, ve := range verrs {
		msgs = append(msgs, fmt.Sprintf("%s: %s", ve.Path, ve.Msg))
	}
	all := strings.Join(msgs, "\n")
	return fmt.Errorf("%d validation errors:\n%s", len(verrs), all)
}

func printSetupVerbose(w io.Writer, doc SetupDocument, relPath string) {
	if doc.GoBlock == nil {
		fmt.Fprintf(w, "%s\n", relPath)
		return
	}
	block := doc.GoBlock
	var parts []string
	if block.Types["Request"] {
		parts = append(parts, "Request")
	}
	if block.Types["Response"] {
		parts = append(parts, "Response")
	}
	if block.Setup != nil {
		parts = append(parts, "Setup")
	}
	if block.Run != nil {
		parts = append(parts, "Run")
	}
	if len(block.Helpers) > 0 {
		parts = append(parts, fmt.Sprintf("%d helpers", len(block.Helpers)))
	}
	fmt.Fprintf(w, "%s — %s", relPath, strings.Join(parts, ", "))
	if block.Run != nil {
		fmt.Fprintf(w, " (defines Run)")
	}
	fmt.Fprintln(w)
}

func runSource(setupDocs []SetupDocument) string {
	for i := len(setupDocs) - 1; i >= 0; i-- {
		doc := setupDocs[i]
		if doc.GoBlock != nil && doc.GoBlock.Run != nil {
			return doc.Path
		}
	}
	return "none"
}

func readSetup(path string) (SetupDocument, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return SetupDocument{Path: path}, nil
		}
		return SetupDocument{}, err
	}
	return ParseSetupDocument(path, string(content))
}

func setupChain(root, leafDir string, doctestDoc SetupDocument) ([]SetupDocument, error) {
	rel, err := filepath.Rel(root, leafDir)
	if err != nil {
		return nil, err
	}
	var parts []string
	if rel != "." {
		parts = strings.Split(rel, string(filepath.Separator))
	}
	var docs []SetupDocument
	if doctestDoc.GoBlock != nil {
		doctestDoc.Path = "DOCTEST.md"
		docs = append(docs, doctestDoc)
	}
	ancestorHelpers := make(map[string]bool)
	for i := 0; i <= len(parts); i++ {
		dir := filepath.Join(append([]string{root}, parts[:i]...)...)
		path := filepath.Join(dir, "SETUP.md")
		doc, err := readSetup(path)
		if err != nil {
			return nil, err
		}
		relPath, _ := filepath.Rel(root, path)
		if doc.GoBlock != nil {
			if v := rules.CheckChildNoRedefine(doc.GoBlock.Types, relPath, i); v != nil {
				return nil, fmt.Errorf("%s: %s", v.Path, v.Msg)
			}
			if v := rules.CheckChildNoRedefineRun(doc.GoBlock.Run != nil, relPath, i); v != nil {
				return nil, fmt.Errorf("%s: %s", v.Path, v.Msg)
			}
			var childHelpers []string
			for _, h := range doc.GoBlock.Helpers {
				childHelpers = append(childHelpers, h.Name)
			}
			if v := rules.CheckNoHelperRedefinition(ancestorHelpers, childHelpers, relPath, i); v != nil {
				return nil, fmt.Errorf("%s: %s", v.Path, v.Msg)
			}
			for _, h := range doc.GoBlock.Helpers {
				ancestorHelpers[h.Name] = true
			}
		}
		doc.Path = relPath
		docs = append(docs, doc)
	}
	return docs, nil
}

func CaseName(path string) string {
	if path == "" {
		return "root"
	}
	return strings.NewReplacer("/", "_", string(filepath.Separator), "_", "-", "_").Replace(filepath.Base(path))
}

func TestFileName(tc TreeCase) string {
	return CaseName(tc.Path) + "_test.go"
}

func TestFuncName(tc TreeCase) string {
	name := CaseName(tc.Path)
	var b strings.Builder
	b.WriteString("TestGeneratedCase")
	upperNext := true
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			if upperNext && r >= 'a' && r <= 'z' {
				r = r - 'a' + 'A'
			}
			b.WriteRune(r)
			upperNext = false
			continue
		}
		upperNext = true
	}
	return b.String()
}

func FindModuleRoot(dir string) (modRoot string, modPath string, ok bool) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", "", false
	}
	for {
		modFile := filepath.Join(dir, "go.mod")
		data, err := os.ReadFile(modFile)
		if err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "module ") {
					modPath = strings.TrimSpace(strings.TrimPrefix(line, "module "))
					return dir, modPath, true
				}
			}
			return "", "", false
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", "", false
		}
		dir = parent
	}
}

// CasesImportInternalPackage reports whether any case imports a path under
// parentModulePath/internal/ using parsed import specs from the setup chain.
func CasesImportInternalPackage(cases []TreeCase, parentModulePath string) bool {
	if parentModulePath == "" {
		return false
	}
	prefix := parentModulePath + "/internal/"
	for _, tc := range cases {
		imports := collectImports(tc.SetupFiles, tc.AssertFile.GoBlock)
		for _, spec := range imports {
			if strings.HasPrefix(spec.Path, prefix) {
				return true
			}
		}
	}
	return false
}

// NewInternalCompileRoot creates a temp compile directory under moduleRoot.
func NewInternalCompileRoot(moduleRoot string) (string, error) {
	absRoot, err := filepath.Abs(moduleRoot)
	if err != nil {
		return "", err
	}
	return os.MkdirTemp(absRoot, ".doctest_run_")
}

// CopyGeneratedTree copies generated files from src to dst for gen-dir review dumps.
// go.mod and go.sum files are not copied.
func CopyGeneratedTree(src, dst string) error {
	absSrc, err := filepath.Abs(src)
	if err != nil {
		return err
	}
	absDst, err := filepath.Abs(dst)
	if err != nil {
		return err
	}
	return filepath.Walk(absSrc, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(absSrc, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		target := filepath.Join(absDst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode().Perm())
		}
		if info.Name() == "go.mod" || info.Name() == "go.sum" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode().Perm())
	})
}

func WriteGoMod(genDir, modRoot, modPath string, hasMod bool) error {
	modFile := filepath.Join(genDir, "go.mod")
	if _, err := os.Stat(modFile); err == nil {
		return nil
	}
	var content string
	if hasMod {
		content = fmt.Sprintf("module testcase\n\ngo 1.21\n\nreplace %s => %s\n", modPath, modRoot)
	} else {
		content = "module testcase\n\ngo 1.21\n"
	}
	if hasMod {
		if extraReplaces := readExtraReplaces(modRoot, modPath); extraReplaces != "" {
			content += extraReplaces
		}
	}
	if err := os.WriteFile(modFile, []byte(content), 0644); err != nil {
		return err
	}
	if hasMod {
		srcGoSum := filepath.Join(modRoot, "go.sum")
		if data, err := os.ReadFile(srcGoSum); err == nil {
			if err := os.WriteFile(filepath.Join(genDir, "go.sum"), data, 0644); err != nil {
				return err
			}
		}
	}
	return nil
}

func readExtraReplaces(modRoot, mainModPath string) string {
	goModPath := filepath.Join(modRoot, "go.mod")
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return ""
	}
	lines := strings.Split(string(data), "\n")
	inReplace := false
	var replaces []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "replace (" {
			inReplace = true
			continue
		}
		if inReplace {
			if trimmed == ")" {
				inReplace = false
				continue
			}
			if strings.HasPrefix(trimmed, "replace ") {
				trimmed = strings.TrimPrefix(trimmed, "replace ")
			}
			replaces = append(replaces, trimmed)
			continue
		}
		if strings.HasPrefix(trimmed, "replace ") {
			replaces = append(replaces, strings.TrimPrefix(trimmed, "replace "))
		}
	}
	var result strings.Builder
	for _, r := range replaces {
		if strings.HasPrefix(r, mainModPath+" ") || strings.HasPrefix(r, mainModPath+"\t") || r == mainModPath {
			continue
		}
		parts := strings.Fields(r)
		if len(parts) >= 3 {
			arrowIdx := len(parts) - 2
			if arrowIdx >= 1 && parts[arrowIdx] == "=>" {
				absPath := parts[len(parts)-1]
				if !filepath.IsAbs(absPath) {
					absPath = filepath.Join(modRoot, absPath)
				}
				absPath = filepath.Clean(absPath)
				modPath := strings.Join(parts[:arrowIdx], " ")
				result.WriteString(fmt.Sprintf("replace %s => %s\n", modPath, absPath))
			}
		}
	}
	return result.String()
}

func TidyGoMod(genDir string) error {
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = genDir
	if out, err := tidy.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod tidy: %v\n%s", err, string(out))
	}
	return nil
}

func CondTidyGoMod(genDir string) error {
	markerFile := filepath.Join(genDir, "doctest.tidy-done")
	if _, err := os.Stat(markerFile); err == nil {
		return nil
	}
	if err := TidyGoMod(genDir); err != nil {
		return err
	}
	f, err := os.Create(markerFile)
	if err != nil {
		return err
	}
	return f.Close()
}

func ResolvePkgUnderTest(root string) (srcDir string, origPkgName string, ok bool) {
	rootSetupPath := filepath.Join(root, "SETUP.md")
	content, err := os.ReadFile(rootSetupPath)
	if err != nil {
		return "", "", false
	}
	pkgName := ""
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- Package under test:") {
			pkgName = strings.Trim(line[len("- Package under test:"):], " `")
		}
	}
	if pkgName == "" {
		return "", "", false
	}
	modRoot, _, hasMod := FindModuleRoot(root)
	if !hasMod {
		return "", "", false
	}
	absDir := filepath.Join(modRoot, pkgName)
	if srcDir, origPkg, ok := readGoFilesForPkg(absDir, pkgName); ok {
		return srcDir, origPkg, true
	}
	entries, err := os.ReadDir(modRoot)
	if err != nil {
		return "", "", false
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dir := filepath.Join(modRoot, entry.Name())
		if srcDir, origPkg, ok := readGoFilesForPkg(dir, pkgName); ok {
			return srcDir, origPkg, true
		}
	}
	return "", "", false
}

func readGoFilesForPkg(dir string, expectedPkg string) (srcDir string, origPkgName string, ok bool) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", "", false
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			continue
		}
		text := string(data)
		var pkgLine string
		if strings.HasPrefix(text, "package ") {
			pkgLine = text[len("package "):]
		} else {
			i := strings.Index(text, "\npackage ")
			if i < 0 {
				continue
			}
			pkgLine = text[i+len("\npackage "):]
		}
		j := strings.IndexAny(pkgLine, " \t\n\r;")
		if j < 0 {
			origPkgName = strings.TrimSpace(pkgLine)
		} else {
			origPkgName = pkgLine[:j]
		}
		if origPkgName == expectedPkg {
			return dir, origPkgName, true
		}
	}
	return "", "", false
}

func CopySourceFiles(genDir, srcDir, origPkgName string) (string, error) {
	if err := os.MkdirAll(genDir, 0755); err != nil {
		return "", err
	}
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return "", err
	}
	newPkgName := origPkgName + "_tc"
	newPkgDecl := "package " + newPkgName
	oldPkgDecl := "package " + origPkgName
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(srcDir, name))
		if err != nil {
			return "", err
		}
		content := strings.Replace(string(data), oldPkgDecl, newPkgDecl, 1)
		dst := filepath.Join(genDir, name)

		existing, _ := os.ReadFile(dst)
		if string(existing) == content {
			continue
		}

		if err := os.WriteFile(dst, []byte(content), 0644); err != nil {
			return "", err
		}
	}
	return newPkgName, nil
}

func GenDirForLeaf(mappingGenRoot, moduleRoot, absLeafDir string) (string, error) {
	rel, err := filepath.Rel(moduleRoot, absLeafDir)
	if err != nil {
		return "", fmt.Errorf("compute relative path for leaf %s from %s: %w", absLeafDir, moduleRoot, err)
	}
	return filepath.Join(mappingGenRoot, rel), nil
}

func MappingGenRoot(absDoctestDir string) (string, string) {
	modRoot, _, _ := FindModuleRoot(absDoctestDir)
	absModRoot, _ := filepath.Abs(absDoctestDir)
	if modRoot != "" {
		absModRoot, _ = filepath.Abs(modRoot)
	}
	return absModRoot, modRoot
}

func WriteGeneratedCase(leafDir string, tc TreeCase, compileOnly bool, pkgName string, docTestRoot string) (string, error) {
	if err := os.MkdirAll(leafDir, 0755); err != nil {
		return "", err
	}
	src, err := AssembleTestSource(tc, compileOnly, pkgName, docTestRoot)
	if err != nil {
		return "", fmt.Errorf("%s: %w", tc.Path, err)
	}
	testFile := TestFileName(tc)
	testPath := filepath.Join(leafDir, testFile)

	tmpFile, err := os.CreateTemp("", ".doctest-gen-*")
	if err != nil {
		return "", err
	}
	tmpPath := tmpFile.Name()
	if _, err := tmpFile.WriteString(src); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return "", err
	}
	tmpFile.Close()

	// Use golang.org/x/tools/imports.Process instead of the goimports binary
	// to avoid external binary dependency.
	// gofmt is not used because imports.Process handles both formatting
	// (via go/format internally, same as gofmt) and import cleanup
	// (adds missing imports, removes unused ones) in a single pass.
	srcBytes, err := os.ReadFile(tmpPath)
	if err != nil {
		os.Remove(tmpPath)
		return "", err
	}
	res, err := imports.Process(tmpPath, srcBytes, nil)
	if err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("format imports failed: %w", err)
	}
	if err := os.WriteFile(tmpPath, res, 0644); err != nil {
		os.Remove(tmpPath)
		return "", err
	}

	existing, _ := os.ReadFile(testPath)
	if string(existing) == string(res) {
		os.Remove(tmpPath)
		return testPath, nil
	}

	if err := os.Rename(tmpPath, testPath); err != nil {
		os.Remove(tmpPath)
		return "", err
	}
	return testPath, nil
}

func CacheMappingGenRoot(absDoctestDir string) (string, string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", "", err
	}
	absModRoot, _ := MappingGenRoot(absDoctestDir)
	mappingRoot := filepath.Join(cacheDir, "doctest", "mapping-gen", absModRoot)
	return mappingRoot, absModRoot, nil
}

func FilterBySubDir(cases []TreeCase, root, subDir string) []TreeCase {
	if subDir == "" || subDir == root {
		return cases
	}
	relSubDir, err := filepath.Rel(root, subDir)
	if err != nil || relSubDir == "." {
		return cases
	}
	prefix := relSubDir + string(filepath.Separator)
	var filtered []TreeCase
	for _, tc := range cases {
		if tc.Path == relSubDir || strings.HasPrefix(tc.Path, prefix) {
			filtered = append(filtered, tc)
		}
	}
	return filtered
}
