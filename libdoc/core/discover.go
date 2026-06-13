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

	rootSetupPath := filepath.Join(root, "SETUP.md")
	rootSetupContent, err := os.ReadFile(rootSetupPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // no SETUP.md at root -> no test cases; caller reports "no tests found"
		}
		return nil, err
	}
	rootSetup, err := ParseSetupDocument(rootSetupPath, string(rootSetupContent))
	if err != nil {
		verrs = append(verrs, ValidationError{Path: "SETUP.md", Msg: err.Error()})
	}
	if rootSetup.GoBlock == nil {
		verrs = append(verrs, ValidationError{Path: "SETUP.md", Msg: "must have a Go code block"})
	} else {
		if v := rules.CheckRootHasRequestResponse(rootSetup.GoBlock.Types, "SETUP.md"); v != nil {
			verrs = append(verrs, ValidationError{Path: v.Path, Msg: v.Msg})
		}
		if v := rules.CheckRootHasSetupOrRun(rootSetup.GoBlock.Setup != nil, rootSetup.GoBlock.Run != nil, "SETUP.md"); v != nil {
			verrs = append(verrs, ValidationError{Path: v.Path, Msg: v.Msg})
		}
	}

	if w != nil {
		printSetupVerbose(w, rootSetup, "SETUP.md")
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
			} else if doc.GoBlock.Setup == nil && doc.GoBlock.Run == nil {
				rel, _ := filepath.Rel(root, setupPath)
				verrs = append(verrs, ValidationError{Path: rel, Msg: "must have func Setup or func Run"})
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
		setupDocs, chainErr := setupChain(root, leafDir)
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

func setupChain(root, leafDir string) ([]SetupDocument, error) {
	rel, err := filepath.Rel(root, leafDir)
	if err != nil {
		return nil, err
	}
	var parts []string
	if rel != "." {
		parts = strings.Split(rel, string(filepath.Separator))
	}
	var docs []SetupDocument
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
	return strings.NewReplacer("/", "_", string(filepath.Separator), "_", "-", "_").Replace(path)
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

func WriteGoMod(genDir, modRoot, modPath string, hasMod bool) error {
	var content string
	if hasMod {
		content = fmt.Sprintf("module testcase\n\ngo 1.21\n\nrequire %s v0.0.0\n\nreplace %s => %s\n", modPath, modPath, modRoot)
	} else {
		content = "module testcase\n\ngo 1.21\n"
	}
	if err := os.WriteFile(filepath.Join(genDir, "go.mod"), []byte(content), 0644); err != nil {
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
	entries, err := os.ReadDir(absDir)
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
		data, err := os.ReadFile(filepath.Join(absDir, name))
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
		if origPkgName != "" {
			return absDir, origPkgName, true
		}
	}
	return "", "", false
}

func CopySourceFiles(genDir, srcDir, origPkgName string) (string, error) {
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
		if err := os.WriteFile(dst, []byte(content), 0644); err != nil {
			return "", err
		}
	}
	return newPkgName, nil
}

func WriteGeneratedCases(dir string, cases []TreeCase, compileOnly bool, w io.Writer, pkgName string, docTestRoot string) ([]string, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	var testPaths []string
	var testFiles []string
	for _, tc := range cases {
		src, err := AssembleTestSource(tc, compileOnly, pkgName, docTestRoot)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", tc.Path, err)
		}
		testFile := TestFileName(tc)
		testPath := filepath.Join(dir, testFile)
		if err := os.WriteFile(testPath, []byte(src), 0644); err != nil {
			return nil, err
		}
		if w != nil {
			fmt.Fprintf(w, "→ %s\n", testPath)
		}
		testPaths = append(testPaths, testPath)
		testFiles = append(testFiles, testFile)
	}
	args := append([]string{"-w"}, testPaths...)
	gofmtCmd := exec.Command("gofmt", args...)
	gofmtCmd.Dir = dir
	if out, err := gofmtCmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("gofmt failed: %v\n%s", err, string(out))
	}
	return testFiles, nil
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
