package core

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/imports"
)

// Ref-mode package layout under gen root (module path is always "testcase"
// when WriteGoMod is used for the external gen tree):
//
//	genRoot/
//	  go.mod
//	  <treeRel>/__droot/droot.go          package droot — types, Run, root helpers
//	  <treeRel>/<parent>/setup.go         intermediate package (non-test)
//	  <treeRel>/<parent>/<mid>/setup.go   nested intermediate package
//	  <treeRel>/<leaf>/…_test.go          thin tests: leaf-local only + import ancestors
//
// treeRel is the path of the doctest root relative to the module root (or "." →
// packages at gen root). This keeps multi-tree ./... + shared GenDir (cold-cache)
// from overwriting each other's packages.
const (
	RefRootDirName = "__droot"
	RefRootPkgName = "droot"
	// RefRootImportPath is the legacy flat import used when treeRel is ".".
	RefRootImportPath = "testcase/__droot"
	// RefIntermediateFileName is the stable non-test filename for intermediate packages.
	RefIntermediateFileName = "setup.go"
)

// RefRootImportForTree returns the go import path for the tree-scoped droot package.
// treeRel is filepath.Rel(modRoot, doctestRoot); empty/"." → testcase/__droot.
func RefRootImportForTree(treeRel string) string {
	treeRel = filepath.ToSlash(filepath.Clean(treeRel))
	if treeRel == "" || treeRel == "." {
		return RefRootImportPath
	}
	return "testcase/" + strings.TrimPrefix(treeRel, "./") + "/" + RefRootDirName
}

// RefRootDirForTree returns the filesystem directory for __droot under genRoot.
func RefRootDirForTree(genRoot, treeRel string) string {
	treeRel = filepath.Clean(treeRel)
	if treeRel == "" || treeRel == "." {
		return filepath.Join(genRoot, RefRootDirName)
	}
	return filepath.Join(genRoot, treeRel, RefRootDirName)
}

// isRootSetupDoc reports whether a setup-chain document belongs to the tree
// root (shared package) rather than a leaf/intermediate directory.
func isRootSetupDoc(doc SetupDocument) bool {
	p := filepath.ToSlash(doc.Path)
	return p == "DOCTEST.md" || p == "SETUP.md" || p == ""
}

// RefIntermediateGroup is one intermediate directory's SETUP docs (shared package).
// Dir is relative to the doctest root (filepath slash form), e.g. "feature/mid".
type RefIntermediateGroup struct {
	Dir  string
	Docs []SetupDocument
}

// RefSetupPartition is the hierarchical split of a case's setup chain for ref mode.
type RefSetupPartition struct {
	RootDocs     []SetupDocument
	Intermediate []RefIntermediateGroup // parents first (root → leaf order)
	LeafDocs     []SetupDocument
}

// SplitRefSetupDocs partitions the setup chain into root vs non-root docs.
// Prefer PartitionRefSetupDocs for hierarchical intermediate packages; this
// helper remains for callers that only need root vs everything-else.
func SplitRefSetupDocs(setupFiles []SetupDocument) (rootDocs, leafDocs []SetupDocument) {
	for _, doc := range setupFiles {
		if isRootSetupDoc(doc) {
			rootDocs = append(rootDocs, doc)
		} else {
			leafDocs = append(leafDocs, doc)
		}
	}
	return rootDocs, leafDocs
}

// PartitionRefSetupDocs splits tc.SetupFiles into root, ordered intermediate
// groups (strict path prefixes of tc.Path), and leaf-local docs (dir == tc.Path).
func PartitionRefSetupDocs(tc TreeCase) RefSetupPartition {
	return PartitionRefSetupDocsFrom(tc.Path, tc.SetupFiles)
}

// PartitionRefSetupDocsFrom is the path-based form of PartitionRefSetupDocs.
func PartitionRefSetupDocsFrom(leafPath string, setupFiles []SetupDocument) RefSetupPartition {
	leafPath = cleanRelDir(leafPath)
	var part RefSetupPartition
	groups := make(map[string][]SetupDocument)
	var order []string
	for _, doc := range setupFiles {
		if isRootSetupDoc(doc) {
			part.RootDocs = append(part.RootDocs, doc)
			continue
		}
		// Docs without a Go block never produce packages or inlined code.
		if doc.GoBlock == nil {
			continue
		}
		dir := setupDocDir(doc)
		if dir == leafPath {
			part.LeafDocs = append(part.LeafDocs, doc)
			continue
		}
		if leafPath != "" && isStrictPathPrefix(dir, leafPath) {
			if _, ok := groups[dir]; !ok {
				order = append(order, dir)
			}
			groups[dir] = append(groups[dir], doc)
			continue
		}
		// Fallback: treat as leaf-local (matches historical non-root-in-leaf behavior).
		part.LeafDocs = append(part.LeafDocs, doc)
	}
	sort.SliceStable(order, func(i, j int) bool {
		di, dj := pathDepth(order[i]), pathDepth(order[j])
		if di != dj {
			return di < dj
		}
		return order[i] < order[j]
	})
	for _, dir := range order {
		part.Intermediate = append(part.Intermediate, RefIntermediateGroup{
			Dir:  dir,
			Docs: groups[dir],
		})
	}
	return part
}

// CollectUniqueRefIntermediates returns unique intermediate groups across cases,
// ordered parents-first so packages can be written with parent imports resolved.
func CollectUniqueRefIntermediates(cases []TreeCase) []RefIntermediateGroup {
	seen := make(map[string]RefIntermediateGroup)
	var order []string
	for _, tc := range cases {
		part := PartitionRefSetupDocs(tc)
		for _, g := range part.Intermediate {
			if _, ok := seen[g.Dir]; ok {
				continue
			}
			seen[g.Dir] = g
			order = append(order, g.Dir)
		}
	}
	sort.SliceStable(order, func(i, j int) bool {
		di, dj := pathDepth(order[i]), pathDepth(order[j])
		if di != dj {
			return di < dj
		}
		return order[i] < order[j]
	})
	out := make([]RefIntermediateGroup, 0, len(order))
	for _, dir := range order {
		out = append(out, seen[dir])
	}
	return out
}

func setupDocDir(doc SetupDocument) string {
	p := filepath.ToSlash(doc.Path)
	if p == "" || p == "DOCTEST.md" || p == "SETUP.md" {
		return ""
	}
	return cleanRelDir(filepath.Dir(p))
}

func cleanRelDir(p string) string {
	p = filepath.ToSlash(filepath.Clean(p))
	if p == "." || p == "/" {
		return ""
	}
	return strings.TrimPrefix(p, "./")
}

func pathDepth(p string) int {
	p = cleanRelDir(p)
	if p == "" {
		return 0
	}
	return strings.Count(p, "/") + 1
}

// isStrictPathPrefix reports whether dir is a strict path prefix of leaf
// (dir != leaf and leaf starts with dir+"/").
func isStrictPathPrefix(dir, leaf string) bool {
	dir = cleanRelDir(dir)
	leaf = cleanRelDir(leaf)
	if dir == "" || leaf == "" || dir == leaf {
		return false
	}
	return strings.HasPrefix(leaf, dir+"/")
}

// parentRelDir returns the parent of dir relative to doctest root, or "" for top-level.
func parentRelDir(dir string) string {
	dir = cleanRelDir(dir)
	if dir == "" {
		return ""
	}
	return cleanRelDir(filepath.Dir(dir))
}

// RefTreeImportPrefix returns the import path prefix for a tree (testcase or testcase/<treeRel>).
func RefTreeImportPrefix(rootImport string) string {
	if rootImport == "" {
		rootImport = RefRootImportPath
	}
	rootImport = filepath.ToSlash(rootImport)
	suffix := "/" + RefRootDirName
	if strings.HasSuffix(rootImport, suffix) {
		return strings.TrimSuffix(rootImport, suffix)
	}
	if rootImport == "testcase/"+RefRootDirName || rootImport == RefRootImportPath {
		return "testcase"
	}
	return "testcase"
}

// RefIntermediateImport returns the go import path for an intermediate package.
// dirRel is relative to the doctest root (e.g. "feature/mid").
func RefIntermediateImport(rootImport, dirRel string) string {
	dirRel = cleanRelDir(dirRel)
	if dirRel == "" {
		return RefTreeImportPrefix(rootImport)
	}
	return RefTreeImportPrefix(rootImport) + "/" + dirRel
}

// RefIntermediateDirForTree returns the filesystem directory for an intermediate package.
func RefIntermediateDirForTree(genRoot, treeRel, dirRel string) string {
	dirRel = filepath.Clean(dirRel)
	base := genRoot
	treeRel = filepath.Clean(treeRel)
	if treeRel != "" && treeRel != "." {
		base = filepath.Join(genRoot, treeRel)
	}
	if dirRel == "" || dirRel == "." {
		return base
	}
	return filepath.Join(base, dirRel)
}

// predeclaredGoIdents are universe identifiers that must not be used as package
// names or import aliases. An alias like `error` shadows the builtin type in
// signatures such as `err error` (seen with intermediate dirs named "error").
var predeclaredGoIdents = map[string]bool{
	"any": true, "bool": true, "byte": true, "comparable": true,
	"complex64": true, "complex128": true, "error": true,
	"float32": true, "float64": true,
	"int": true, "int8": true, "int16": true, "int32": true, "int64": true,
	"rune": true, "string": true,
	"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true, "uintptr": true,
	"true": true, "false": true, "iota": true, "nil": true,
	"append": true, "cap": true, "clear": true, "close": true, "complex": true, "copy": true,
	"delete": true, "imag": true, "len": true, "make": true, "max": true, "min": true,
	"new": true, "panic": true, "print": true, "println": true, "real": true, "recover": true,
}

// SanitizePackageName turns a directory basename into a valid Go package identifier.
func SanitizePackageName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "pkg"
	}
	var b strings.Builder
	for i, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			if i == 0 || b.Len() == 0 {
				b.WriteByte('p')
			}
			b.WriteRune(r)
		case r == '_':
			if b.Len() == 0 {
				b.WriteByte('p')
			}
			b.WriteByte('_')
		default:
			if b.Len() > 0 {
				b.WriteByte('_')
			}
		}
	}
	out := b.String()
	out = strings.Trim(out, "_")
	if out == "" {
		return "pkg"
	}
	// Leading digit after trim should not happen; keep safe.
	if out[0] >= '0' && out[0] <= '9' {
		out = "p" + out
	}
	if token.IsKeyword(out) || predeclaredGoIdents[out] {
		return "pkg_" + out
	}
	// Avoid colliding with the reserved root package name when nested oddly.
	if out == RefRootPkgName || out == RefRootDirName {
		return "pkg_" + out
	}
	return out
}

// RefIntermediatePkgName is the package clause name for an intermediate dir (basename).
func RefIntermediatePkgName(dirRel string) string {
	dirRel = cleanRelDir(dirRel)
	base := filepath.Base(dirRel)
	if base == "" || base == "." {
		return "pkg"
	}
	return SanitizePackageName(base)
}

// RefIntermediateAlias is a unique import alias for an intermediate dir (path-based).
func RefIntermediateAlias(dirRel string) string {
	dirRel = cleanRelDir(dirRel)
	if dirRel == "" {
		return "pkg"
	}
	return SanitizePackageName(strings.ReplaceAll(dirRel, "/", "_"))
}

// intermediateSetupName returns the exported Setup symbol for the i-th setup in a group.
func intermediateSetupName(i, total int) string {
	if total <= 1 {
		return "Setup"
	}
	return fmt.Sprintf("Setup%d", i)
}

// resolveRefParent picks the nearest ancestor intermediate package, or droot.
// intermediateByDir maps dir → group for packages that will exist on disk.
func resolveRefParent(dir string, intermediateByDir map[string]RefIntermediateGroup, rootImport string) (parentImport, parentAlias string, parentDocs []SetupDocument, parentIsRoot bool) {
	chain := resolveRefAncestorChain(dir, intermediateByDir)
	if len(chain) == 0 {
		if rootImport == "" {
			rootImport = RefRootImportPath
		}
		return rootImport, RefRootPkgName, nil, true
	}
	nearest := chain[len(chain)-1]
	return RefIntermediateImport(rootImport, nearest.Dir), RefIntermediateAlias(nearest.Dir), nearest.Docs, false
}

// resolveRefAncestorChain returns intermediate ancestor groups for dir, ordered
// parents-first (rootward → nearest parent). Empty means only droot is parent.
func resolveRefAncestorChain(dir string, intermediateByDir map[string]RefIntermediateGroup) []RefIntermediateGroup {
	var up []RefIntermediateGroup
	p := parentRelDir(dir)
	for p != "" {
		if g, ok := intermediateByDir[p]; ok {
			up = append(up, g)
		}
		p = parentRelDir(p)
	}
	// up is nearest-first; reverse to parents-first.
	for i, j := 0, len(up)-1; i < j; i, j = i+1, j-1 {
		up[i], up[j] = up[j], up[i]
	}
	return up
}

// collectRootSymbolRenames maps unexported root package symbols to exported
// names so leaf packages can reference them as droot.ExportedName.
// Exported symbols map to themselves. Types Request/Response are excluded
// (already qualified separately as droot.Request).
func collectRootSymbolRenames(rootDocs []SetupDocument) map[string]string {
	renames := make(map[string]string)
	add := func(name string) {
		if name == "" || name == "_" ||
			name == "Request" || name == "Response" || name == "Setup" || name == "Run" || name == "Assert" {
			return
		}
		// Never rewrite blank identifiers or tiny builtins via symbol map.
		if name == "err" || name == "ok" || name == "t" || name == "req" || name == "resp" {
			return
		}
		if token.IsExported(name) {
			renames[name] = name
			return
		}
		renames[name] = exportIdent(name)
	}
	for _, doc := range rootDocs {
		if doc.GoBlock == nil {
			continue
		}
		for _, h := range doc.GoBlock.Helpers {
			add(h.Name)
		}
		// Methods stay package-local on root types; still export name if leaf needed (rare).
		for _, m := range doc.GoBlock.Methods {
			// Method name only; receiver type stays in droot.
			add(m.Name)
		}
		for _, decl := range doc.GoBlock.VarDecls {
			for _, n := range declIdentNames(decl) {
				add(n)
			}
		}
		for _, decl := range doc.GoBlock.Consts {
			for _, n := range declIdentNames(decl) {
				add(n)
			}
		}
	}
	return renames
}

func exportIdent(name string) string {
	if name == "" || token.IsExported(name) {
		return name
	}
	r, size := utf8.DecodeRuneInString(name)
	return string(unicode.ToUpper(r)) + name[size:]
}

func declIdentNames(decl string) []string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", "package p\n"+decl, 0)
	if err != nil {
		return nil
	}
	var names []string
	for _, d := range f.Decls {
		gd, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, id := range vs.Names {
				names = append(names, id.Name)
			}
		}
	}
	return names
}

// rewriteBareIdents renames bare identifiers according to renames (old→new).
// Uses go/scanner so string/rune/comment contents are left untouched (critical:
// fixture generators embed Go source in string literals).
func rewriteBareIdents(src string, renames map[string]string) string {
	if len(renames) == 0 || src == "" {
		return src
	}
	return scanReplaceIdents(src, func(name string, afterDot bool) string {
		if afterDot {
			return name
		}
		if neu, ok := renames[name]; ok {
			return neu
		}
		return name
	})
}

// qualifyRootSymbols prefixes bare root symbols with alias.ExportedName
// (e.g. runType → droot.RunType). Skips string/comment content.
func qualifyRootSymbols(src, alias string, renames map[string]string) string {
	if (len(renames) == 0 && alias == "") || src == "" {
		return src
	}
	// Map both bare and already-exported forms to the exported form.
	lookup := map[string]string{}
	for old, neu := range renames {
		lookup[old] = neu
		lookup[neu] = neu
	}
	return scanReplaceIdents(src, func(name string, afterDot bool) string {
		if afterDot || name == "_" || name == "" {
			return name
		}
		exported, ok := lookup[name]
		if !ok {
			return name
		}
		return alias + "." + exported
	})
}

// scanReplaceIdents rewrites IDENT tokens via repl, skipping strings/comments.
// afterDot is true when the previous non-comment token was '.'.
// Token spans are taken from consecutive scan offsets so operator tokens
// (empty lit) cannot desync the rewrite from the source.
func scanReplaceIdents(src string, repl func(name string, afterDot bool) string) string {
	type tok struct {
		off int
		tok token.Token
		lit string
	}
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	var s scanner.Scanner
	s.Init(file, []byte(src), nil, scanner.ScanComments)
	var toks []tok
	for {
		pos, t, lit := s.Scan()
		if t == token.EOF {
			break
		}
		// Skip automatically inserted semicolons — keep original newlines via gaps.
		if t == token.SEMICOLON && lit == "\n" {
			continue
		}
		toks = append(toks, tok{off: fset.File(pos).Offset(pos), tok: t, lit: lit})
	}
	var out strings.Builder
	out.Grow(len(src) + 64)
	prevDot := false
	lastEnd := 0
	for i, tk := range toks {
		end := len(src)
		if i+1 < len(toks) {
			end = toks[i+1].off
		}
		// Prefer source slice for non-IDENT so spacing stays exact.
		if tk.off > lastEnd {
			out.WriteString(src[lastEnd:tk.off])
		}
		switch tk.tok {
		case token.IDENT:
			// Token text length is len(lit); do not use end-off (includes trailing space).
			neu := repl(tk.lit, prevDot)
			out.WriteString(neu)
			lastEnd = tk.off + len(tk.lit)
			prevDot = false
		case token.PERIOD:
			out.WriteByte('.')
			lastEnd = tk.off + 1
			prevDot = true
		default:
			// Copy exact source bytes for this token (and only the token, not trailing space).
			tokEnd := tk.off
			if tk.lit != "" {
				tokEnd = tk.off + len(tk.lit)
			} else {
				// Operators: length of token.String(), e.g. ":=", "..."
				tokEnd = tk.off + len(tk.tok.String())
				if tokEnd > end {
					tokEnd = end
				}
				// Prefer raw source slice when it matches.
				if tokEnd <= len(src) {
					out.WriteString(src[tk.off:tokEnd])
					lastEnd = tokEnd
					prevDot = false
					continue
				}
			}
			if tokEnd > len(src) {
				tokEnd = len(src)
			}
			out.WriteString(src[tk.off:tokEnd])
			lastEnd = tokEnd
			prevDot = false
		}
	}
	if lastEnd < len(src) {
		out.WriteString(src[lastEnd:])
	}
	return out.String()
}

// AssembleRefRootSource emits the shared root package: types, consts/vars,
// helpers/methods, and package-level Run from root setup documents.
// Unexported package-level symbols are exported so leaf packages can use them.
func AssembleRefRootSource(rootDocs []SetupDocument, pkgName string) (string, error) {
	if pkgName == "" {
		pkgName = RefRootPkgName
	}
	var run *FuncSnippet
	for _, doc := range rootDocs {
		if doc.GoBlock != nil && doc.GoBlock.Run != nil {
			runCopy := *doc.GoBlock.Run
			run = &runCopy
			break
		}
	}
	if run == nil {
		return "", fmt.Errorf("missing Run(t *testing.T, req *Request) (*Response, error) in root setup chain")
	}

	renames := collectRootSymbolRenames(rootDocs)

	var buf strings.Builder
	buf.WriteString("package ")
	buf.WriteString(pkgName)
	buf.WriteString("\n\n")

	importsMap := collectImports(rootDocs, GoBlock{})
	if _, ok := importsMap["testing"]; !ok {
		importsMap["testing"] = &ImportSpec{Path: "testing"}
	}
	if _, ok := importsMap[sessionImportPath]; !ok {
		importsMap[sessionImportPath] = &ImportSpec{Path: sessionImportPath}
	}
	writeImportBlock(&buf, importsMap)

	// Session context is injected as d *session.Doctest on Setup/Run — no package free DOCTEST_* vars.

	// Emit types/methods/vars/helpers then rewrite unexported symbols to exported.
	var body strings.Builder
	writePackageLevelTypesAndMethods(&body, rootDocs, GoBlock{})
	writePackageLevelConstVars(&body, rootDocs, GoBlock{})
	writePackageLevelHelpers(&body, rootDocs, GoBlock{})

	for i, doc := range rootDocs {
		if doc.GoBlock == nil || doc.GoBlock.Setup == nil {
			continue
		}
		fn := *doc.GoBlock.Setup
		if fn.Name == "" || fn.Name == "Setup" {
			fn.Name = fmt.Sprintf("RootSetup%d", i)
		}
		fn.Params = ensureDoctestParam(fn.Params)
		writePackageFunc(&body, fn)
		body.WriteString("\n")
	}

	run.Name = "Run"
	run.Params = ensureDoctestParam(run.Params)
	writePackageFunc(&body, *run)
	body.WriteString("\n")

	// Export unexported package-level symbols (vars/helpers) and rewrite refs.
	bodyStr := rewriteBareIdents(body.String(), renames)
	buf.WriteString(bodyStr)
	return buf.String(), nil
}

// AssembleRefIntermediateSource emits a non-test intermediate package from
// SETUP docs for one directory. It imports droot and every intermediate
// ancestor (parents-first), exports Setup + helpers (same export rename pattern
// as root), and qualifies ancestor/root symbols so multi-level chains can call
// grandparent helpers (e.g. cold-cache → gen-dir-mode → auto).
//
// ancestors is ordered parents-first (rootward → nearest parent); empty means
// only droot is the parent.
func AssembleRefIntermediateSource(
	docs []SetupDocument,
	pkgName string,
	rootImport, rootAlias string,
	rootDocs []SetupDocument,
	ancestors []RefIntermediateGroup,
) (string, error) {
	if pkgName == "" {
		pkgName = "pkg"
	}
	if rootImport == "" {
		rootImport = RefRootImportPath
	}
	if rootAlias == "" {
		rootAlias = RefRootPkgName
	}

	rootRenames := collectRootSymbolRenames(rootDocs)
	rootTypes := rootTypeNamesForQualify(rootDocs)
	ownRenames := collectRootSymbolRenames(docs)

	var buf strings.Builder
	buf.WriteString("package ")
	buf.WriteString(pkgName)
	buf.WriteString("\n\n")

	importsMap := collectImports(docs, GoBlock{})
	if _, ok := importsMap["testing"]; !ok {
		importsMap["testing"] = &ImportSpec{Path: "testing"}
	}
	if _, ok := importsMap[sessionImportPath]; !ok {
		importsMap[sessionImportPath] = &ImportSpec{Path: sessionImportPath}
	}
	importsMap[rootAlias+"\x00"+rootImport] = &ImportSpec{Name: rootAlias, Path: rootImport}
	for _, g := range ancestors {
		alias := RefIntermediateAlias(g.Dir)
		imp := RefIntermediateImport(rootImport, g.Dir)
		importsMap[alias+"\x00"+imp] = &ImportSpec{Name: alias, Path: imp}
	}
	writeImportBlock(&buf, importsMap)

	var body strings.Builder
	writePackageLevelTypesAndMethods(&body, docs, GoBlock{})
	writePackageLevelConstVars(&body, docs, GoBlock{})
	writePackageLevelHelpers(&body, docs, GoBlock{})

	setupTotal := 0
	for _, doc := range docs {
		if doc.GoBlock != nil && doc.GoBlock.Setup != nil {
			setupTotal++
		}
	}
	setupIdx := 0
	for _, doc := range docs {
		if doc.GoBlock == nil || doc.GoBlock.Setup == nil {
			continue
		}
		fn := *doc.GoBlock.Setup
		fn.Name = intermediateSetupName(setupIdx, setupTotal)
		setupIdx++
		fn.Params = ensureDoctestParam(fn.Params)
		writePackageFunc(&body, fn)
		body.WriteString("\n")
	}

	bodyStr := rewriteBareIdents(body.String(), ownRenames)
	// Qualify ancestor symbols parents-first (same order as leaf qualifyAncestorSymbols).
	// Only types actually declared on each intermediate — never default
	// Request/Response (those live on droot).
	for _, g := range ancestors {
		alias := RefIntermediateAlias(g.Dir)
		bodyStr = qualifyRootSymbols(bodyStr, alias, collectRootSymbolRenames(g.Docs))
		bodyStr = qualifyRootTypes(bodyStr, alias, collectRootTypeNames(g.Docs))
	}
	bodyStr = qualifyRootSymbols(bodyStr, rootAlias, rootRenames)
	bodyStr = qualifyRootTypes(bodyStr, rootAlias, rootTypes)
	buf.WriteString(bodyStr)
	return buf.String(), nil
}

// qualifyAncestorSymbols rewrites bare symbols/types from each intermediate
// ancestor and from root into alias-qualified form (for leaf-local code).
func qualifyAncestorSymbols(src string, part RefSetupPartition, rootAlias string, rootRenames map[string]string, rootTypes []string) string {
	if src == "" {
		return src
	}
	// Parents first is fine; helper names are unique across the chain.
	for _, g := range part.Intermediate {
		alias := RefIntermediateAlias(g.Dir)
		src = qualifyRootSymbols(src, alias, collectRootSymbolRenames(g.Docs))
		src = qualifyRootTypes(src, alias, collectRootTypeNames(g.Docs))
	}
	src = qualifyRootSymbols(src, rootAlias, rootRenames)
	src = qualifyRootTypes(src, rootAlias, rootTypes)
	return src
}

// AssembleRefLeafTestSource emits a thin leaf *_test.go that imports the root
// package and intermediate ancestor packages. Leaf must not redefine root types
// or inline intermediate SETUP bodies.
func AssembleRefLeafTestSource(tc TreeCase, compileOnly bool, pkgName, docTestRoot, rootImport, rootAlias string) (string, error) {
	if pkgName == "" {
		pkgName = "testcase"
	}
	if rootImport == "" {
		rootImport = RefRootImportPath
	}
	if rootAlias == "" {
		rootAlias = RefRootPkgName
	}

	part := PartitionRefSetupDocs(tc)
	rootDocs := part.RootDocs
	leafDocs := part.LeafDocs
	renames := collectRootSymbolRenames(rootDocs)
	rootTypes := rootTypeNamesForQualify(rootDocs)

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
	writeImportBlock(&buf, importsMap)

	// Leaf-only types/helpers — rewrite references to root + intermediate symbols.
	var leafBlob strings.Builder
	writePackageLevelTypesAndMethods(&leafBlob, leafDocs, tc.AssertFile.GoBlock)
	writePackageLevelConstVars(&leafBlob, leafDocs, tc.AssertFile.GoBlock)
	writePackageLevelHelpers(&leafBlob, leafDocs, tc.AssertFile.GoBlock)
	leafTop := qualifyAncestorSymbols(leafBlob.String(), part, rootAlias, renames, rootTypes)
	buf.WriteString(leafTop)

	buf.WriteString("func ")
	buf.WriteString(TestFuncName(tc))
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
		fn.Params = qualifyAncestorSymbols(fn.Params, part, rootAlias, renames, rootTypes)
		fn.Results = qualifyAncestorSymbols(fn.Results, part, rootAlias, renames, rootTypes)
		fn.ResultTypes = qualifyAncestorSymbols(fn.ResultTypes, part, rootAlias, renames, rootTypes)
		fn.ClosureResults = qualifyAncestorSymbols(fn.ClosureResults, part, rootAlias, renames, rootTypes)
		fn.Body = qualifyAncestorSymbols(fn.Body, part, rootAlias, renames, rootTypes)
		fn.Params = ensureDoctestParam(fn.Params)
		writeFuncClosure(&buf, name, fn)
		buf.WriteString(fmt.Sprintf("\tif err := %s(t, d, req); err != nil {\n", name))
		buf.WriteString(fmt.Sprintf("\t\tt.Fatalf(\"%s failed: %%v\", err)\n", escapeString(doc.Path)))
		buf.WriteString("\t}\n")
	}

	assertFn := *tc.AssertFile.GoBlock.Assert
	assertFn.Params = qualifyAncestorSymbols(assertFn.Params, part, rootAlias, renames, rootTypes)
	assertFn.Results = qualifyAncestorSymbols(assertFn.Results, part, rootAlias, renames, rootTypes)
	assertFn.ResultTypes = qualifyAncestorSymbols(assertFn.ResultTypes, part, rootAlias, renames, rootTypes)
	assertFn.ClosureResults = qualifyAncestorSymbols(assertFn.ClosureResults, part, rootAlias, renames, rootTypes)
	assertFn.Body = qualifyAncestorSymbols(assertFn.Body, part, rootAlias, renames, rootTypes)
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

func writeImportBlock(buf *strings.Builder, imports map[string]*ImportSpec) {
	if len(imports) == 0 {
		return
	}
	buf.WriteString("import (\n")
	importList := make([]*ImportSpec, 0, len(imports))
	for _, spec := range imports {
		importList = append(importList, spec)
	}
	sort.Slice(importList, func(i, j int) bool {
		if importList[i].Path != importList[j].Path {
			return importList[i].Path < importList[j].Path
		}
		return importList[i].Name < importList[j].Name
	})
	for _, spec := range importList {
		if spec.Name != "" {
			buf.WriteString("\t" + spec.Name + " \"" + spec.Path + "\"\n")
		} else {
			buf.WriteString("\t\"" + spec.Path + "\"\n")
		}
	}
	buf.WriteString(")\n\n")
}

// collectRootTypeNames returns type names declared in root docs (Request, Response, WarnCase, …).
func collectRootTypeNames(rootDocs []SetupDocument) []string {
	seen := map[string]bool{}
	var names []string
	for _, doc := range rootDocs {
		if doc.GoBlock == nil {
			continue
		}
		for name := range doc.GoBlock.Types {
			if name == "" || seen[name] {
				continue
			}
			seen[name] = true
			names = append(names, name)
		}
	}
	sort.Slice(names, func(i, j int) bool { return len(names[i]) > len(names[j]) })
	return names
}

// qualifyRootTypes rewrites bare type identifiers in rootTypeNames to alias.X.
// An empty rootTypeNames means no rewrites (callers must pass Request/Response
// explicitly for root — do not default here, or intermediate packages with no
// local types would steal Request/Response onto the wrong alias).
// Skips string/comment content so fixture generators keep raw `*Request` text.
func qualifyRootTypes(s, rootAlias string, rootTypeNames []string) string {
	if s == "" || len(rootTypeNames) == 0 {
		return s
	}
	typeSet := map[string]bool{}
	for _, t := range rootTypeNames {
		typeSet[t] = true
	}
	return scanReplaceIdents(s, func(name string, afterDot bool) string {
		if afterDot || !typeSet[name] {
			return name
		}
		return rootAlias + "." + name
	})
}

// rootTypeNamesForQualify returns type names declared in root docs, always
// including Request and Response so leaf/intermediate Setup params qualify.
func rootTypeNamesForQualify(rootDocs []SetupDocument) []string {
	names := collectRootTypeNames(rootDocs)
	seen := map[string]bool{}
	for _, n := range names {
		seen[n] = true
	}
	for _, req := range []string{"Request", "Response"} {
		if !seen[req] {
			names = append(names, req)
			seen[req] = true
		}
	}
	// Prefer longer names first (stable with collectRootTypeNames).
	sort.Slice(names, func(i, j int) bool { return len(names[i]) > len(names[j]) })
	return names
}

func qualifyRootTypesInBody(body, rootAlias string, rootTypeNames []string) string {
	return qualifyRootTypes(body, rootAlias, rootTypeNames)
}

// WriteRefTree writes the shared root package, intermediate packages once, and
// thin leaf tests for each case. Flat layout (treeRel ".") under genRoot.
func WriteRefTree(genRoot string, cases []TreeCase, docTestRoot string, compileOnly bool, pkgName string) error {
	if len(cases) == 0 {
		return fmt.Errorf("WriteRefTree: no cases")
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
	// Flat layout without mod/tree scoping: __droot at gen root.
	rootDir := filepath.Join(genRoot, RefRootDirName)
	rootImport := RefRootImportPath
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		return err
	}
	rootPath := filepath.Join(rootDir, "droot.go")
	if err := WriteFormattedGo(rootPath, rootSrc); err != nil {
		return fmt.Errorf("write ref root package: %w", err)
	}

	if err := WriteRefIntermediatePackages(genRoot, ".", rootImport, rootDocs, cases); err != nil {
		return err
	}

	for _, tc := range cases {
		leafDir := genRoot
		if tc.Path != "" {
			leafDir = filepath.Join(genRoot, tc.Path)
		}
		if _, err := WriteRefLeafCase(leafDir, tc, compileOnly, pkgName, docTestRoot, rootImport); err != nil {
			return err
		}
	}
	return nil
}

// WriteRefIntermediatePackages writes each unique intermediate package once under
// genRoot (tree-scoped when treeRel is not ".").
func WriteRefIntermediatePackages(genRoot, treeRel, rootImport string, rootDocs []SetupDocument, cases []TreeCase) error {
	if rootImport == "" {
		rootImport = RefRootImportPath
	}
	groups := CollectUniqueRefIntermediates(cases)
	byDir := make(map[string]RefIntermediateGroup, len(groups))
	for _, g := range groups {
		byDir[g.Dir] = g
	}
	for _, g := range groups {
		ancestors := resolveRefAncestorChain(g.Dir, byDir)
		pkgName := RefIntermediatePkgName(g.Dir)
		src, err := AssembleRefIntermediateSource(
			g.Docs, pkgName,
			rootImport, RefRootPkgName, rootDocs,
			ancestors,
		)
		if err != nil {
			return fmt.Errorf("assemble intermediate %s: %w", g.Dir, err)
		}
		dir := RefIntermediateDirForTree(genRoot, treeRel, g.Dir)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		path := filepath.Join(dir, RefIntermediateFileName)
		if err := WriteFormattedGo(path, src); err != nil {
			return fmt.Errorf("write intermediate package %s: %w", g.Dir, err)
		}
	}
	return nil
}

// WriteRefLeafCase writes a thin leaf test file under leafDir.
func WriteRefLeafCase(leafDir string, tc TreeCase, compileOnly bool, pkgName, docTestRoot, rootImport string) (string, error) {
	if err := os.MkdirAll(leafDir, 0755); err != nil {
		return "", err
	}
	if rootImport == "" {
		rootImport = RefRootImportPath
	}
	src, err := AssembleRefLeafTestSource(tc, compileOnly, pkgName, docTestRoot, rootImport, RefRootPkgName)
	if err != nil {
		return "", fmt.Errorf("%s: %w", tc.Path, err)
	}
	testPath := filepath.Join(leafDir, TestFileName(tc))
	if err := WriteFormattedGo(testPath, src); err != nil {
		return "", err
	}
	return testPath, nil
}

// WriteFormattedGo formats src with imports.Process and writes atomically when changed.
func WriteFormattedGo(path, src string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	res, err := imports.Process(path, []byte(src), nil)
	if err != nil {
		_ = os.WriteFile(path, []byte(src), 0644)
		return fmt.Errorf("format imports for %s: %w", path, err)
	}
	existing, _ := os.ReadFile(path)
	if string(existing) == string(res) {
		return nil
	}
	tmpFile, err := os.CreateTemp(filepath.Dir(path), ".doctest-gen-*")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	if _, err := tmpFile.Write(res); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return err
	}
	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}
	if err := os.Rename(tmpPath, path); err != nil {
		if !isCrossDeviceRename(err) {
			os.Remove(tmpPath)
			return err
		}
		if writeErr := os.WriteFile(path, res, 0644); writeErr != nil {
			os.Remove(tmpPath)
			return fmt.Errorf("rename (cross-device): %w; write fallback: %v", err, writeErr)
		}
		os.Remove(tmpPath)
	}
	return nil
}


