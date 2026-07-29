package core

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"syscall"

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
	// Memoize SETUP parse by abs path within this walk (shared ancestors).
	setupCache := make(map[string]setupCacheEntry)

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
	if _, rootStatErr := os.Stat(rootSetupPath); rootStatErr == nil {
		rootSetup, rootSetupErr := readSetupCached(rootSetupPath, setupCache)
		if rootSetupErr != nil {
			verrs = append(verrs, ValidationError{Path: "SETUP.md", Msg: rootSetupErr.Error()})
		} else if rootSetup.GoBlock != nil {
			if v := rules.CheckRootSetupNoRequestResponseRun(rootSetup.GoBlock.Types, rootSetup.GoBlock.Run != nil, "SETUP.md"); v != nil {
				verrs = append(verrs, ValidationError{Path: v.Path, Msg: v.Msg})
			}
		}
		if w != nil && rootSetupErr == nil {
			printSetupVerbose(w, rootSetup, "SETUP.md")
		}
	} else if !os.IsNotExist(rootStatErr) {
		return nil, rootStatErr
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
					doc, readErr := readSetupCached(setupPath, setupCache)
					if readErr == nil && doc.GoBlock != nil {
						printedSetupDirs[setupPath] = true
						printSetupVerbose(w, doc, filepath.Join(relPath, "SETUP.md"))
					}
				}
			}
			// Intermediate dirs: missing SETUP.md is OK (pure nested-tree parents,
			// empty grouping dirs). When SETUP.md exists, require Go Setup.
			setupPath := filepath.Join(path, "SETUP.md")
			if _, statErr := os.Stat(setupPath); os.IsNotExist(statErr) {
				return nil
			} else if statErr != nil {
				rel, _ := filepath.Rel(root, setupPath)
				verrs = append(verrs, ValidationError{Path: rel, Msg: statErr.Error()})
				return nil
			}
			doc, readErr := readSetupCached(setupPath, setupCache)
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
		setupDocs, chainErr := setupChainCached(root, leafDir, doctestDoc, setupCache)
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
			Name:        CaseName(relLeaf),
			Path:        relLeaf,
			SetupFiles:  setupDocs,
			AssertFile:  assertDoc,
			Labels:      append([]string(nil), assertDoc.Frontmatter.Labels...),
			Explanation: assertDoc.Frontmatter.Explanation,
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

type setupCacheEntry struct {
	doc SetupDocument
	err error
}

func readSetup(path string) (SetupDocument, error) {
	return readSetupCached(path, nil)
}

// readSetupCached loads and parses SETUP.md once per path when cache != nil.
// Missing file: empty document, nil error (same as historical readSetup).
func readSetupCached(path string, cache map[string]setupCacheEntry) (SetupDocument, error) {
	key := path
	if abs, aerr := filepath.Abs(path); aerr == nil {
		key = abs
	}
	if cache != nil {
		if e, ok := cache[key]; ok {
			return e.doc, e.err
		}
	}
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			doc := SetupDocument{Path: path}
			if cache != nil {
				cache[key] = setupCacheEntry{doc: doc}
			}
			return doc, nil
		}
		if cache != nil {
			cache[key] = setupCacheEntry{err: err}
		}
		return SetupDocument{}, err
	}
	doc, err := ParseSetupDocument(path, string(content))
	if cache != nil {
		cache[key] = setupCacheEntry{doc: doc, err: err}
	}
	return doc, err
}

func setupChain(root, leafDir string, doctestDoc SetupDocument) ([]SetupDocument, error) {
	return setupChainCached(root, leafDir, doctestDoc, nil)
}

func setupChainCached(root, leafDir string, doctestDoc SetupDocument, cache map[string]setupCacheEntry) ([]SetupDocument, error) {
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
		doc, err := readSetupCached(path, cache)
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

// genModMu serializes WriteGoMod / CondTidyGoMod / gen-manifest flush so
// parallel ./... trees that share one gen root (e.g. --cold-cache
// mapping-gen-cold) do not race on go.mod / go.sum / tidy markers / manifest.
var genModMu sync.Mutex

// writeFileIfChanged writes data only when content differs, preserving mtime
// when unchanged (critical for go test package cache inputs under gen root).
func writeFileIfChanged(path string, data []byte, perm os.FileMode) (wrote bool, err error) {
	existing, err := os.ReadFile(path)
	if err == nil && string(existing) == string(data) {
		return false, nil
	}
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}
	if err := os.WriteFile(path, data, perm); err != nil {
		return false, err
	}
	return true, nil
}

// WriteGoMod builds nested go.mod / go.sum under genDir. Skip key is the
// content hash of the desired final bytes in doctest.gen-manifest (not
// doctest.gomod-fp). tidy-done is invalidated only when go.mod/go.sum actually wrote.
// Assert/session replaces are omitted for the doctest self-module so multi-tree
// ./... prepare with differing ineffective flags does not churn mtimes.
func WriteGoMod(genDir, modRoot, modPath string, hasMod bool, withAssertReplace bool, assertCacheDir string, withSessionReplace bool, sessionCacheDir string) error {
	genModMu.Lock()
	defer genModMu.Unlock()

	genDir = absGenRoot(genDir)
	if err := os.MkdirAll(genDir, 0755); err != nil {
		return err
	}

	// Prefer the source module's go directive so tidy under kool with-go1.19
	// (or any older toolchain) does not fail with "go.mod indicates go 1.21".
	goLine := "go 1.19"
	if hasMod {
		if v := readGoDirective(modRoot); v != "" {
			goLine = "go " + v
		}
	}
	var content string
	if hasMod {
		content = fmt.Sprintf("module testcase\n\n%s\n\nreplace %s => %s\n", goLine, modPath, modRoot)
	} else {
		content = fmt.Sprintf("module testcase\n\n%s\n", goLine)
	}
	// Modules parent already replaces (path or module→module): vendor inject
	// must not dual-replace the same path (parent wins), except private-fork
	// offline safety where vendor FS replace of A suppresses parent module→module.
	var extraReplaces string
	var parentPathReplaced, parentModuleReplaced map[string]bool
	if hasMod {
		extraReplaces, parentPathReplaced, parentModuleReplaced = readExtraReplaces(modRoot, modPath)
	}
	// Effective-only: raw flags that do not affect content (doctest self-module)
	// must not produce different desired bytes — multi-tree ./... prepare calls
	// WriteGoMod with different per-tree flags.
	if withAssertReplace && assertCacheDir != "" && modPath != "github.com/xhd2015/doctest" {
		content += fmt.Sprintf("replace %s => %s\n", AssertImportPath, assertCacheDir)
	}
	if withSessionReplace && sessionCacheDir != "" && modPath != "github.com/xhd2015/doctest" {
		content += fmt.Sprintf("replace %s => %s\n", SessionImportPath, sessionCacheDir)
	}

	// When modRoot has vendor/modules.txt, inject xgo-style require+replace for
	// each vendored module and ensure placeholder go.mod under modules that
	// lack one. Always on when vendor/ exists — no user flag. Skip vendor
	// replace for modules already covered by parent path replace or (non-fork)
	// parent module→module replace.
	parentGoVer := strings.TrimPrefix(goLine, "go ")
	vendorExtra, suppressParentModule, err := vendorBridgeForModRoot(modRoot, parentGoVer, parentPathReplaced, parentModuleReplaced)
	if err != nil {
		return err
	}
	// Parent replaces first (minus any suppressed by private-fork vendor FS win),
	// then vendor require/replace so offline vendor targets are last writer for
	// those paths when they override parent module→module.
	if extraReplaces != "" {
		content += filterReplaceLinesByLHS(extraReplaces, suppressParentModule)
	}
	if vendorExtra != "" {
		content += vendorExtra
	}

	man, err := cachedGenManifestLocked(genDir)
	if err != nil {
		return err
	}

	modWrote, err := man.writeRelIfChanged(genDir, "go.mod", []byte(content))
	if err != nil {
		return err
	}

	sumWrote := false
	if hasMod {
		srcGoSum := filepath.Join(modRoot, "go.sum")
		if data, err := os.ReadFile(srcGoSum); err == nil {
			w, werr := man.writeRelIfChanged(genDir, "go.sum", data)
			if werr != nil {
				return werr
			}
			sumWrote = w
		} else if os.IsNotExist(err) {
			sumPath := filepath.Join(genDir, "go.sum")
			if _, serr := os.Stat(sumPath); serr == nil {
				if rerr := os.Remove(sumPath); rerr != nil && !os.IsNotExist(rerr) {
					return rerr
				}
				sumWrote = true
			}
			man.deleteHash("go.sum")
		} else {
			return err
		}
	}

	if err := man.flush(genDir); err != nil {
		return err
	}
	if modWrote || sumWrote {
		os.Remove(filepath.Join(genDir, "doctest.tidy-done"))
	}
	return nil
}

// readExtraReplaces copies parent go.mod replace directives into gen form
// (absolute-izing filesystem path targets). Returns:
//   - extra replace text for gen go.mod
//   - parentPathReplaced: left-hand paths with filesystem path RHS
//   - parentModuleReplaced: left-hand paths with module→module RHS
// Vendor inject uses these so parent path replace always wins, parent
// module→module wins unless private-fork offline safety prefers vendor FS.
func readExtraReplaces(modRoot, mainModPath string) (string, map[string]bool, map[string]bool) {
	parentPathReplaced := make(map[string]bool)
	parentModuleReplaced := make(map[string]bool)
	goModPath := filepath.Join(modRoot, "go.mod")
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return "", parentPathReplaced, parentModuleReplaced
	}
	lines := strings.Split(string(data), "\n")
	inReplace := false
	var replaces []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}
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
		// Skip comments / empty after trim
		if r == "" || strings.HasPrefix(r, "//") {
			continue
		}
		// Drop replace of the main module itself (WriteGoMod already emits that).
		if strings.HasPrefix(r, mainModPath+" ") || strings.HasPrefix(r, mainModPath+"\t") || r == mainModPath {
			continue
		}
		parts := strings.Fields(r)
		arrowIdx := -1
		for i, p := range parts {
			if p == "=>" {
				arrowIdx = i
				break
			}
		}
		if arrowIdx < 1 || arrowIdx+1 >= len(parts) {
			continue
		}
		// Left may be "path" or "path version"; module path is the first token.
		leftModPath := parts[0]
		left := strings.Join(parts[:arrowIdx], " ")
		rightParts := parts[arrowIdx+1:]
		// Path replace: single token that looks like a filesystem path.
		// Module→module replace: "module/path v1.2.3" (keep as-is).
		if len(rightParts) == 1 {
			right := rightParts[0]
			if isFilesystemReplaceTarget(right) {
				if !filepath.IsAbs(right) {
					right = filepath.Join(modRoot, right)
				}
				right = filepath.Clean(right)
				result.WriteString(fmt.Sprintf("replace %s => %s\n", left, right))
				if leftModPath != "" {
					parentPathReplaced[leftModPath] = true
				}
				continue
			}
			// Single-token module path without version — still valid replace.
			result.WriteString(fmt.Sprintf("replace %s => %s\n", left, right))
			if leftModPath != "" {
				parentModuleReplaced[leftModPath] = true
			}
			continue
		}
		// module => other/module vX.Y.Z  (or more tokens; preserve verbatim)
		result.WriteString(fmt.Sprintf("replace %s => %s\n", left, strings.Join(rightParts, " ")))
		if leftModPath != "" {
			parentModuleReplaced[leftModPath] = true
		}
	}
	return result.String(), parentPathReplaced, parentModuleReplaced
}

// filterReplaceLinesByLHS drops "replace …" lines whose left-hand module path
// (first field after "replace ") is in drop. Used when private-fork vendor FS
// replace of A supersedes parent module→module for A.
func filterReplaceLinesByLHS(replaces string, drop map[string]bool) string {
	if replaces == "" || len(drop) == 0 {
		return replaces
	}
	var b strings.Builder
	for _, line := range strings.Split(replaces, "\n") {
		if line == "" {
			continue
		}
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "replace ") {
			rest := strings.TrimSpace(strings.TrimPrefix(trimmed, "replace "))
			fields := strings.Fields(rest)
			if len(fields) >= 1 && drop[fields[0]] {
				continue
			}
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}

// readGoDirective returns the version from the first "go X.Y" line in modRoot/go.mod.
func readGoDirective(modRoot string) string {
	data, err := os.ReadFile(filepath.Join(modRoot, "go.mod"))
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) >= 2 && fields[0] == "go" {
			return fields[1]
		}
	}
	return ""
}

// isFilesystemReplaceTarget reports whether right-hand side of a replace is a
// local path (absolutize against modRoot) rather than a module path.
func isFilesystemReplaceTarget(right string) bool {
	if right == "" {
		return false
	}
	if filepath.IsAbs(right) {
		return true
	}
	if strings.HasPrefix(right, "./") || strings.HasPrefix(right, "../") {
		return true
	}
	// Bare relative dir/file occasionally used in replace (rare).
	if strings.HasPrefix(right, ".") {
		return true
	}
	return false
}

// TidyGoMod runs `go mod tidy` in genDir. When goCache is non-empty it is
// applied as GOCACHE via ChildEnv (key-replace) so cold-cache isolation does
// not require process os.Setenv.
func TidyGoMod(genDir string, goCache string) error {
	// Caller must hold genModMu when concurrent trees share genDir.
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = genDir
	if goCache != "" {
		tidy.Env = ChildEnv(nil, "GOCACHE="+goCache)
	}
	if out, err := tidy.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod tidy: %v\n%s", err, string(out))
	}
	return nil
}

// CondTidyGoMod runs TidyGoMod once per genDir (marker file), serializing via
// genModMu. goCache is optional isolated GOCACHE for the tidy child.
func CondTidyGoMod(genDir string, goCache string) error {
	genModMu.Lock()
	defer genModMu.Unlock()

	markerFile := filepath.Join(genDir, "doctest.tidy-done")
	if _, err := os.Stat(markerFile); err == nil {
		// Marker already present: bookkeeping unchanged this run.
		NoteDesired(genDir, "doctest.tidy-done")
		NoteWrite(genDir, "doctest.tidy-done", false)
		noteGenBookkeeping(genDir)
		return nil
	}
	if err := TidyGoMod(genDir, goCache); err != nil {
		return err
	}
	f, err := os.Create(markerFile)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	NoteDesired(genDir, "doctest.tidy-done")
	NoteWriteEx(genDir, "doctest.tidy-done", true, true) // newly created marker
	noteGenBookkeeping(genDir)
	return nil
}

// noteGenBookkeeping marks go.mod / go.sum / tidy-done / manifest as desired
// when present (tidy may rewrite go.sum outside writeRelIfChanged).
// Paths without a prior write outcome default to unchanged (already on disk).
func noteGenBookkeeping(genDir string) {
	for _, name := range []string{"go.mod", "go.sum", "doctest.tidy-done", genManifestFile} {
		if _, err := os.Stat(filepath.Join(genDir, name)); err == nil {
			NoteDesired(genDir, name)
			// Only set unchanged when no stronger outcome was recorded.
			NoteWrite(genDir, name, false)
		}
	}
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

		_, stErr := os.Stat(dst)
		existed := stErr == nil
		existing, _ := os.ReadFile(dst)
		wrote := false
		if string(existing) != content {
			if err := os.WriteFile(dst, []byte(content), 0644); err != nil {
				return "", err
			}
			wrote = true
		}
		// Desired even on content-identical skip.
		if root, rel, ok := findGenRootWithManifest(dst); ok {
			NoteDesired(root, rel)
			NoteWriteEx(root, rel, wrote, wrote && !existed)
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

// WriteGeneratedCase writes the generated test file for one leaf.
// wrote is true when the on-disk content changed (or the file was created).
func WriteGeneratedCase(leafDir string, tc TreeCase, compileOnly bool, pkgName string, docTestRoot string) (path string, wrote bool, err error) {
	if err := os.MkdirAll(leafDir, 0755); err != nil {
		return "", false, err
	}
	src, err := AssembleTestSource(tc, compileOnly, pkgName, docTestRoot)
	if err != nil {
		return "", false, fmt.Errorf("%s: %w", tc.Path, err)
	}
	testFile := TestFileName(tc)
	testPath := filepath.Join(leafDir, testFile)

	// Format in memory. Do NOT stage a temp file inside leafDir first:
	// CreateTemp+Remove there updates the package directory mtime even when
	// content is unchanged, which busts go test's result cache (testlog
	// hashes chdir/stat mtime of dirs under the module root).
	res, err := formatGeneratedGo(testPath, []byte(src))
	if err != nil {
		return "", false, fmt.Errorf("format generated Go failed: %w", err)
	}

	existing, _ := os.ReadFile(testPath)
	if string(existing) == string(res) {
		// Unchanged: leave leafDir completely untouched so GOCACHE can hit.
		return testPath, false, nil
	}

	// Content differs — write atomically via same-dir temp + rename so we
	// avoid EXDEV when /tmp and the gen cache are on different mounts
	// (e.g. GitHub Actions). Fall back to WriteFile only on EXDEV.
	tmpFile, err := os.CreateTemp(leafDir, ".doctest-gen-*")
	if err != nil {
		return "", false, err
	}
	tmpPath := tmpFile.Name()
	if _, err := tmpFile.Write(res); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return "", false, err
	}
	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpPath)
		return "", false, err
	}

	if err := os.Rename(tmpPath, testPath); err != nil {
		// Only fall back to write+remove for EXDEV (cross-device). Other rename
		// failures (permission, missing path, etc.) should surface as-is.
		if !isCrossDeviceRename(err) {
			os.Remove(tmpPath)
			return "", false, err
		}
		if writeErr := os.WriteFile(testPath, res, 0644); writeErr != nil {
			os.Remove(tmpPath)
			return "", false, fmt.Errorf("rename (cross-device): %w; write fallback: %v", err, writeErr)
		}
		os.Remove(tmpPath)
	}
	return testPath, true, nil
}

// isCrossDeviceRename reports whether err is an EXDEV from os.Rename
// ("invalid cross-device link").
func isCrossDeviceRename(err error) bool {
	var linkErr *os.LinkError
	if errors.As(err, &linkErr) {
		return linkErr.Err == syscall.EXDEV
	}
	return errors.Is(err, syscall.EXDEV)
}

func CacheMappingGenRoot(absDoctestDir string) (string, string, error) {
	cacheDir, err := CacheHome()
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
