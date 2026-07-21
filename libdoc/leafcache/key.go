package leafcache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/parser"
	"go/token"
	"hash"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/mod/modfile"
)

// ComputeLeafKey walks the leaf spine, local import closure, module identity,
// and local replace modules to produce a stable lowercase hex digest.
//
// Hashed: AlgoVersion, GoVersion, absolute cleaned TreeRoot (tree identity so
// identical relative content under different roots cannot share a key), module
// go.mod (+ go.sum if present), spine Go blocks only (root DOCTEST, ancestor
// SETUPs on the path to the leaf, leaf SETUP, leaf ASSERT — sibling-branch
// SETUPs are not mixed in), local package sources in the import closure under
// ModuleRoot and local replace modules, and those local replace modules'
// go.mod/go.sum.
//
// Not hashed: remote module source trees (they contribute only via
// go.mod/go.sum identity); process env values (os.Getenv/os.LookupEnv are not
// special-cased); non-spine SETUP.md under the tree.
func ComputeLeafKey(in KeyInput) (string, error) {
	if err := validateKeyInput(in); err != nil {
		return "", err
	}

	h := sha256.New()
	writeField(h, "algo", AlgoVersion)
	writeField(h, "goVersion", in.GoVersion)

	modRoot, err := filepath.Abs(in.ModuleRoot)
	if err != nil {
		return "", fmt.Errorf("leafcache: ModuleRoot: %w", err)
	}
	treeRoot, err := filepath.Abs(in.TreeRoot)
	if err != nil {
		return "", fmt.Errorf("leafcache: TreeRoot: %w", err)
	}
	leafDir, err := filepath.Abs(in.LeafDir)
	if err != nil {
		return "", fmt.Errorf("leafcache: LeafDir: %w", err)
	}

	// Tree identity: absolute cleaned TreeRoot distinguishes twin trees that
	// share the same relative leaf path and content (cross-tree cache isolation).
	writeField(h, "treeRoot", filepath.Clean(treeRoot))

	// Primary module meta (optional: fixtures without go.mod still hash spine).
	localMods := map[string]*moduleInfo{}
	primaryMod, err := loadModule(modRoot)
	if err != nil {
		if !isMissingGoMod(err) {
			return "", err
		}
		// No go.mod: empty module identity; no local package closure under a module.
		writeField(h, "module:path", "")
		writeField(h, "module:gomod", "")
		writeField(h, "module:gosum", "")
	} else {
		writeModuleMeta(h, "module", primaryMod)

		// Local replace modules (BFS): their go.mod/go.sum identity always mixes in.
		localMods, err = collectLocalModules(primaryMod)
		if err != nil {
			return "", err
		}
		// Stable order by module path (skip primary — already written as "module").
		paths := make([]string, 0, len(localMods))
		for p, m := range localMods {
			if m.Dir == primaryMod.Dir {
				continue
			}
			paths = append(paths, p)
		}
		sort.Strings(paths)
		for _, p := range paths {
			writeModuleMeta(h, "replace-module:"+p, localMods[p])
		}
	}

	// Spine Go blocks only (no tree-wide SETUP walk, no osenv value mixing).
	// Trim surrounding whitespace so prose-only SETUP rewrites that only differ
	// by a trailing newline inside the fence stay key-stable.
	spineCodes, err := collectSpine(treeRoot, leafDir)
	if err != nil {
		return "", err
	}
	for _, sb := range spineCodes {
		writeField(h, "spine:"+sb.rel, strings.TrimSpace(sb.code))
	}

	// Local import closure from spine (current module + local replace pkgs).
	seedImports := make(map[string]struct{})
	for _, sb := range spineCodes {
		for _, imp := range extractImports(sb.code) {
			seedImports[imp] = struct{}{}
		}
	}
	pkgPaths, err := localImportClosure(seedImports, localMods)
	if err != nil {
		return "", err
	}
	sort.Strings(pkgPaths)
	for _, pkgPath := range pkgPaths {
		dir, ok := resolveLocalPackage(pkgPath, localMods)
		if !ok {
			continue
		}
		files, err := listPackageGoFiles(dir)
		if err != nil {
			return "", fmt.Errorf("leafcache: package %s: %w", pkgPath, err)
		}
		writeField(h, "pkg", pkgPath)
		for _, f := range files {
			content, err := os.ReadFile(f.path)
			if err != nil {
				return "", err
			}
			writeField(h, "pkgfile:"+pkgPath+"/"+f.name, string(content))
		}
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func validateKeyInput(in KeyInput) error {
	if strings.TrimSpace(in.ModuleRoot) == "" {
		return fmt.Errorf("leafcache: ModuleRoot is required")
	}
	if strings.TrimSpace(in.TreeRoot) == "" {
		return fmt.Errorf("leafcache: TreeRoot is required")
	}
	if strings.TrimSpace(in.LeafDir) == "" {
		return fmt.Errorf("leafcache: LeafDir is required")
	}
	return nil
}

func writeField(h hash.Hash, name, value string) {
	// Length-prefixed fields avoid ambiguity between adjacent concatenations.
	fmt.Fprintf(h, "%d:%s\n%d:%s\n", len(name), name, len(value), value)
}

// moduleInfo is one go module root (primary or local replace).
type moduleInfo struct {
	Path    string // module path from go.mod
	Dir     string // absolute directory
	GoMod   []byte
	GoSum   []byte // nil if absent
	Replace []localReplace
}

type localReplace struct {
	OldPath string
	NewDir  string // absolute
}

func isMissingGoMod(err error) bool {
	if err == nil {
		return false
	}
	if os.IsNotExist(err) {
		return true
	}
	// Wrapped read errors from loadModule.
	return strings.Contains(err.Error(), "go.mod") &&
		(strings.Contains(err.Error(), "no such file") || strings.Contains(err.Error(), "cannot find"))
}

func loadModule(dir string) (*moduleInfo, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	goModPath := filepath.Join(dir, "go.mod")
	data, err := os.ReadFile(goModPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		return nil, fmt.Errorf("leafcache: read %s: %w", goModPath, err)
	}
	f, err := modfile.Parse(goModPath, data, nil)
	if err != nil {
		return nil, fmt.Errorf("leafcache: parse %s: %w", goModPath, err)
	}
	if f.Module == nil || f.Module.Mod.Path == "" {
		return nil, fmt.Errorf("leafcache: %s: missing module path", goModPath)
	}
	info := &moduleInfo{
		Path:  f.Module.Mod.Path,
		Dir:   dir,
		GoMod: data,
	}
	sumPath := filepath.Join(dir, "go.sum")
	if sum, err := os.ReadFile(sumPath); err == nil {
		info.GoSum = sum
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	for _, r := range f.Replace {
		if r == nil {
			continue
		}
		newPath := r.New.Path
		if !isLocalReplacePath(newPath) {
			continue
		}
		abs := newPath
		if !filepath.IsAbs(abs) {
			abs = filepath.Join(dir, filepath.FromSlash(newPath))
		}
		abs, err = filepath.Abs(abs)
		if err != nil {
			return nil, err
		}
		info.Replace = append(info.Replace, localReplace{
			OldPath: r.Old.Path,
			NewDir:  abs,
		})
	}
	return info, nil
}

// isLocalReplacePath reports whether a replace target is local (relative or absolute path).
// Remote versioned replaces use module paths without leading ./ ../ or /.
func isLocalReplacePath(p string) bool {
	if p == "" {
		return false
	}
	if filepath.IsAbs(p) {
		return true
	}
	// Relative: ./foo, ../foo, foo/bar (Go allows bare relative for replace).
	// Module paths look like "example.com/lib" — treat as local only if they
	// start with . or contain only path-like segments that begin with ./.
	if strings.HasPrefix(p, "./") || strings.HasPrefix(p, "../") || p == "." || p == ".." {
		return true
	}
	// Windows drive paths handled by IsAbs; Unix abs by IsAbs.
	// Bare "lib" relative replace is rare; Go docs require ./ or ../ for relative.
	// Still accept paths starting with '.' .
	if strings.HasPrefix(p, ".") {
		return true
	}
	return false
}

// collectLocalModules BFS-walks local replace edges starting from primary.
// Map key is module path from go.mod.
func collectLocalModules(primary *moduleInfo) (map[string]*moduleInfo, error) {
	out := map[string]*moduleInfo{primary.Path: primary}
	queue := []*moduleInfo{primary}
	seenDir := map[string]struct{}{primary.Dir: {}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, r := range cur.Replace {
			if _, ok := seenDir[r.NewDir]; ok {
				continue
			}
			seenDir[r.NewDir] = struct{}{}
			child, err := loadModule(r.NewDir)
			if err != nil {
				// Local replace target without go.mod: still record dir under old path.
				// Skip if unreadable — callers that import it will fail later.
				if os.IsNotExist(err) || strings.Contains(err.Error(), "no such file") {
					continue
				}
				// loadModule wraps errors; try to continue only for missing go.mod.
				if _, statErr := os.Stat(filepath.Join(r.NewDir, "go.mod")); os.IsNotExist(statErr) {
					continue
				}
				return nil, err
			}
			// Index by module path; also keep replace old path for resolution.
			if existing, ok := out[child.Path]; ok {
				if existing.Dir != child.Dir {
					// Prefer first; ignore conflict for key purposes.
				}
			} else {
				out[child.Path] = child
			}
			// Also index under replace old path so imports of the replaced module resolve.
			if _, ok := out[r.OldPath]; !ok {
				// Use a shallow alias pointing at same module info if paths differ.
				out[r.OldPath] = child
			}
			queue = append(queue, child)
		}
	}
	return out, nil
}

func writeModuleMeta(h hash.Hash, tag string, m *moduleInfo) {
	writeField(h, tag+":path", m.Path)
	writeField(h, tag+":gomod", string(m.GoMod))
	if m.GoSum != nil {
		writeField(h, tag+":gosum", string(m.GoSum))
	} else {
		writeField(h, tag+":gosum", "")
	}
}

type spineBlock struct {
	rel  string
	code string
}

// collectSpine returns Go blocks for: root DOCTEST, each SETUP from tree root
// down to leaf (inclusive), and leaf ASSERT. Order is fixed for stability.
func collectSpine(treeRoot, leafDir string) ([]spineBlock, error) {
	var out []spineBlock

	// Root DOCTEST.md (final go block).
	doctestPath := filepath.Join(treeRoot, "DOCTEST.md")
	code, err := extractFinalGoBlockFile(doctestPath)
	if err != nil {
		return nil, fmt.Errorf("leafcache: spine DOCTEST: %w", err)
	}
	out = append(out, spineBlock{rel: "DOCTEST.md", code: code})

	// Ancestor SETUP.md from tree root down to leaf (inclusive).
	relChain, err := relPathChain(treeRoot, leafDir)
	if err != nil {
		return nil, err
	}
	// relChain: ["", "group", "group/leaf"] when leaf is tree/group/leaf
	for _, rel := range relChain {
		dir := treeRoot
		if rel != "" {
			dir = filepath.Join(treeRoot, filepath.FromSlash(rel))
		}
		setupPath := filepath.Join(dir, "SETUP.md")
		if _, statErr := os.Stat(setupPath); os.IsNotExist(statErr) {
			continue
		}
		code, err := extractFinalGoBlockFile(setupPath)
		if err != nil {
			// SETUP without a go block is allowed (empty contribution).
			if isMissingGoBlock(err) {
				continue
			}
			return nil, fmt.Errorf("leafcache: spine %s: %w", setupPath, err)
		}
		label := "SETUP.md"
		if rel != "" {
			label = filepath.ToSlash(filepath.Join(rel, "SETUP.md"))
		}
		out = append(out, spineBlock{rel: label, code: code})
	}

	assertPath := filepath.Join(leafDir, "ASSERT.md")
	code, err = extractFinalGoBlockFile(assertPath)
	if err != nil {
		return nil, fmt.Errorf("leafcache: spine ASSERT: %w", err)
	}
	leafRel, err := filepath.Rel(treeRoot, leafDir)
	if err != nil {
		leafRel = "leaf"
	}
	out = append(out, spineBlock{
		rel:  filepath.ToSlash(filepath.Join(leafRel, "ASSERT.md")),
		code: code,
	})
	return out, nil
}

// relPathChain returns path segments from root to leaf as slash-paths,
// starting with "" (the root itself), then each intermediate dir.
func relPathChain(root, leaf string) ([]string, error) {
	root = filepath.Clean(root)
	leaf = filepath.Clean(leaf)
	rel, err := filepath.Rel(root, leaf)
	if err != nil {
		return nil, fmt.Errorf("leafcache: LeafDir not under TreeRoot: %w", err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return nil, fmt.Errorf("leafcache: LeafDir %q is not under TreeRoot %q", leaf, root)
	}
	out := []string{""}
	if rel == "." {
		return out, nil
	}
	parts := strings.Split(filepath.ToSlash(rel), "/")
	acc := ""
	for _, p := range parts {
		if p == "" || p == "." {
			continue
		}
		if acc == "" {
			acc = p
		} else {
			acc = acc + "/" + p
		}
		out = append(out, acc)
	}
	return out, nil
}

func extractFinalGoBlockFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return extractFinalGoBlock(string(data))
}

func isMissingGoBlock(err error) bool {
	return err != nil && strings.Contains(err.Error(), "missing go block")
}

// extractFinalGoBlock returns the last ```go fenced block's body.
func extractFinalGoBlock(content string) (string, error) {
	blocks := findGoBlocks(content)
	if len(blocks) == 0 {
		return "", fmt.Errorf("missing go block")
	}
	return blocks[len(blocks)-1], nil
}

func findGoBlocks(content string) []string {
	var blocks []string
	i := 0
	for {
		start := strings.Index(content[i:], "```go")
		if start < 0 {
			return blocks
		}
		start += i
		lineEnd := strings.IndexByte(content[start:], '\n')
		if lineEnd < 0 {
			return blocks
		}
		codeStart := start + lineEnd + 1
		close := strings.Index(content[codeStart:], "```")
		if close < 0 {
			return blocks
		}
		close += codeStart
		blocks = append(blocks, content[codeStart:close])
		i = close + len("```")
	}
}

// extractImports parses a fragment that looks like a Go file body (imports + decls)
// and returns import paths. Fragments from doctest often omit package clauses.
func extractImports(code string) []string {
	candidates := []string{
		code,
		"package leafcache_spine\n\n" + code,
	}
	var imports []string
	seen := map[string]struct{}{}
	fset := token.NewFileSet()
	for _, src := range candidates {
		f, err := parser.ParseFile(fset, "spine.go", src, parser.ImportsOnly|parser.ParseComments)
		if err != nil {
			// Try full parse (some fragments need it).
			f, err = parser.ParseFile(fset, "spine.go", src, parser.ParseComments)
			if err != nil {
				continue
			}
		}
		for _, imp := range f.Imports {
			p := strings.Trim(imp.Path.Value, `"`)
			if p == "" {
				continue
			}
			if _, ok := seen[p]; ok {
				continue
			}
			seen[p] = struct{}{}
			imports = append(imports, p)
		}
		if len(imports) > 0 || f != nil {
			break
		}
	}
	return imports
}

// localImportClosure BFS-expands seed import paths to local packages only.
func localImportClosure(seed map[string]struct{}, localMods map[string]*moduleInfo) ([]string, error) {
	var queue []string
	for p := range seed {
		queue = append(queue, p)
	}
	sort.Strings(queue)

	seen := map[string]struct{}{}
	var out []string

	for len(queue) > 0 {
		imp := queue[0]
		queue = queue[1:]
		if _, ok := seen[imp]; ok {
			continue
		}
		dir, ok := resolveLocalPackage(imp, localMods)
		if !ok {
			// stdlib / remote — ignore source tree
			continue
		}
		seen[imp] = struct{}{}
		out = append(out, imp)

		files, err := listPackageGoFiles(dir)
		if err != nil {
			return nil, err
		}
		next := map[string]struct{}{}
		for _, f := range files {
			content, err := os.ReadFile(f.path)
			if err != nil {
				return nil, err
			}
			for _, child := range extractImportsFromFile(string(content)) {
				if _, ok := seen[child]; ok {
					continue
				}
				if _, isLocal := resolveLocalPackage(child, localMods); isLocal {
					next[child] = struct{}{}
				}
			}
		}
		var add []string
		for p := range next {
			add = append(add, p)
		}
		sort.Strings(add)
		queue = append(queue, add...)
	}
	return out, nil
}

func extractImportsFromFile(src string) []string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "pkg.go", src, parser.ImportsOnly)
	if err != nil {
		return nil
	}
	var out []string
	for _, imp := range f.Imports {
		p := strings.Trim(imp.Path.Value, `"`)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// resolveLocalPackage maps an import path to a local directory if it belongs
// to the primary module or a local replace module.
func resolveLocalPackage(importPath string, localMods map[string]*moduleInfo) (string, bool) {
	// Prefer longest matching module path.
	var bestPath string
	var best *moduleInfo
	for modPath, m := range localMods {
		if importPath == modPath || strings.HasPrefix(importPath, modPath+"/") {
			if len(modPath) > len(bestPath) {
				bestPath = modPath
				best = m
			}
		}
	}
	if best == nil {
		return "", false
	}
	suffix := strings.TrimPrefix(importPath, bestPath)
	suffix = strings.TrimPrefix(suffix, "/")
	dir := best.Dir
	if suffix != "" {
		dir = filepath.Join(best.Dir, filepath.FromSlash(suffix))
	}
	// Package dir must exist.
	if st, err := os.Stat(dir); err != nil || !st.IsDir() {
		return "", false
	}
	return dir, true
}

type goFile struct {
	name string
	path string
}

// listPackageGoFiles lists non-test .go files in dir (package sources only).
func listPackageGoFiles(dir string) ([]goFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []goFile
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		// Skip build-ignored files only if we can cheaply detect; include all
		// non-test .go for content sensitivity (//go:build is still content).
		files = append(files, goFile{name: name, path: filepath.Join(dir, name)})
	}
	sort.Slice(files, func(i, j int) bool { return files[i].name < files[j].name })
	return files, nil
}

