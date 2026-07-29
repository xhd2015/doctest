package core

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// XgoRuntimeMockImport is the import path that requires the xgo test driver.
const XgoRuntimeMockImport = "github.com/xhd2015/xgo/runtime/mock"

// ResolveGoTestCmd maps --go-cmd mode and mock detection to a binary name.
// Returns "go" or "xgo" (never a PATH). Empty mode and "auto" use needsXgo.
func ResolveGoTestCmd(mode string, needsXgo bool) (string, error) {
	switch strings.TrimSpace(mode) {
	case "", "auto":
		if needsXgo {
			return "xgo", nil
		}
		return "go", nil
	case "xgo":
		return "xgo", nil
	case "go":
		return "go", nil
	default:
		return "", fmt.Errorf("invalid --go-cmd value %q: must be one of auto, xgo, go", mode)
	}
}

// EnsureGoTestCmdAvailable looks up cmd on searchPATH (or the process PATH when
// searchPATH is empty). When searchPATH is non-empty it is used exclusively —
// no fallback to the process PATH (so tests can inject a fake PATH without xgo).
func EnsureGoTestCmdAvailable(cmd, searchPATH string) error {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return fmt.Errorf("go test command is empty")
	}
	if lookPathIn(cmd, searchPATH) != "" {
		return nil
	}
	if searchPATH != "" {
		return fmt.Errorf("%s not found in PATH", cmd)
	}
	return fmt.Errorf("%s not found in PATH", cmd)
}

// lookPathIn finds an executable named cmd. When searchPATH is non-empty, only
// those directories are searched; otherwise os.Getenv("PATH") is used.
func lookPathIn(cmd, searchPATH string) string {
	var dirs []string
	if searchPATH != "" {
		dirs = filepath.SplitList(searchPATH)
	} else {
		dirs = filepath.SplitList(os.Getenv("PATH"))
	}
	for _, dir := range dirs {
		if dir == "" {
			continue
		}
		candidate := filepath.Join(dir, cmd)
		if isExecutable(candidate) {
			return candidate
		}
	}
	return ""
}

func isExecutable(path string) bool {
	st, err := os.Stat(path)
	if err != nil || st.IsDir() {
		return false
	}
	// Unix: any execute bit. On Windows Stat alone is enough if the file exists.
	return st.Mode()&0o111 != 0 || filepath.Ext(path) != ""
}

// DetectXgoMockUsage reports whether github.com/xhd2015/xgo/runtime/mock is
// imported transitively from entryImportPaths within the module at modRoot.
//
// Only packages under the project module are walked (source under modRoot).
// External imports are not descended into, except that the mock path itself
// (or a subpath) is recognized immediately when seen as an import.
//
// This is a source-level walk (go/parser) so fixtures that blank-import the
// mock path without a real module download still detect correctly.
func DetectXgoMockUsage(modRoot string, entryImportPaths []string) (bool, error) {
	if modRoot == "" {
		return false, fmt.Errorf("DetectXgoMockUsage: empty modRoot")
	}
	modRoot, err := filepath.Abs(modRoot)
	if err != nil {
		return false, err
	}
	_, modPath, ok := FindModuleRoot(modRoot)
	if !ok {
		// Treat modRoot as the module root even without a readable go.mod if
		// the directory exists; without modPath we can only match the mock
		// path literally on the entry list.
		if st, stErr := os.Stat(modRoot); stErr != nil || !st.IsDir() {
			return false, fmt.Errorf("DetectXgoMockUsage: modRoot %s: %w", modRoot, stErr)
		}
	}

	// Queue of import paths to visit.
	queue := make([]string, 0, len(entryImportPaths))
	seen := make(map[string]bool, len(entryImportPaths))
	for _, p := range entryImportPaths {
		p = strings.TrimSpace(p)
		if p == "" || seen[p] {
			continue
		}
		seen[p] = true
		queue = append(queue, p)
	}

	for len(queue) > 0 {
		imp := queue[0]
		queue = queue[1:]

		if isXgoMockImport(imp) {
			return true, nil
		}

		// Only walk packages that live under the project module.
		pkgDir, inMod := packageDirInModule(modRoot, modPath, imp)
		if !inMod {
			continue
		}
		imports, err := listPackageImports(pkgDir)
		if err != nil {
			// Missing package dir: skip (entry may be optional / not materialised).
			if os.IsNotExist(err) {
				continue
			}
			return false, err
		}
		for _, child := range imports {
			if seen[child] {
				continue
			}
			seen[child] = true
			if isXgoMockImport(child) {
				return true, nil
			}
			// Enqueue project packages for further walk; external non-mock
			// imports stay in seen so we do not re-check them.
			if modPath != "" && (child == modPath || strings.HasPrefix(child, modPath+"/")) {
				queue = append(queue, child)
			}
		}
	}
	return false, nil
}

func isXgoMockImport(path string) bool {
	return path == XgoRuntimeMockImport || strings.HasPrefix(path, XgoRuntimeMockImport+"/")
}

// packageDirInModule maps an import path to a directory under modRoot when the
// path belongs to the module. ok is false for external modules.
func packageDirInModule(modRoot, modPath, importPath string) (dir string, ok bool) {
	if modPath == "" {
		return "", false
	}
	if importPath == modPath {
		return modRoot, true
	}
	if !strings.HasPrefix(importPath, modPath+"/") {
		return "", false
	}
	rel := strings.TrimPrefix(importPath, modPath+"/")
	return filepath.Join(modRoot, filepath.FromSlash(rel)), true
}

// listPackageImports parses all .go files in pkgDir (including _test.go) and
// returns unique import paths. Does not follow the package graph.
func listPackageImports(pkgDir string) ([]string, error) {
	entries, err := os.ReadDir(pkgDir)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	seen := make(map[string]bool)
	var out []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		// Skip build-ignored files only if we can cheaply detect //go:build ignore
		// via parse; parser.ParseFile with ImportsOnly is enough for paths.
		path := filepath.Join(pkgDir, name)
		f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			// Unreadable / invalid source: skip file, do not fail the whole walk.
			continue
		}
		for _, is := range f.Imports {
			p := strings.Trim(is.Path.Value, `"`)
			if p == "" || seen[p] {
				continue
			}
			seen[p] = true
			out = append(out, p)
		}
	}
	return out, nil
}

// EntryImportPathsFromCases collects unique import paths from Setup/Assert Go
// blocks of the given cases. Used as detection entrypoints for --go-cmd=auto.
func EntryImportPathsFromCases(cases []TreeCase) []string {
	seen := make(map[string]bool)
	var out []string
	for _, tc := range cases {
		imports := collectImports(tc.SetupFiles, tc.AssertFile.GoBlock)
		for _, spec := range imports {
			if spec.Path == "" || seen[spec.Path] {
				continue
			}
			seen[spec.Path] = true
			out = append(out, spec.Path)
		}
	}
	return out
}

// ResolveAndEnsureGoTestCmd detects (when mode is auto), resolves the binary
// name, and ensures it is available on PATH. modRoot may be empty when mode is
// force go/xgo (detection skipped).
func ResolveAndEnsureGoTestCmd(mode, modRoot string, entryImportPaths []string) (cmd string, needsXgo bool, err error) {
	mode = strings.TrimSpace(mode)
	if mode == "" {
		mode = "auto"
	}
	needsXgo = false
	if mode == "auto" {
		if modRoot != "" && len(entryImportPaths) > 0 {
			needsXgo, err = DetectXgoMockUsage(modRoot, entryImportPaths)
			if err != nil {
				return "", false, err
			}
		}
	}
	cmd, err = ResolveGoTestCmd(mode, needsXgo)
	if err != nil {
		return "", needsXgo, err
	}
	if err := EnsureGoTestCmdAvailable(cmd, ""); err != nil {
		return cmd, needsXgo, err
	}
	return cmd, needsXgo, nil
}

// ValidateGoCmdMode returns an error if mode is not empty, auto, xgo, or go.
func ValidateGoCmdMode(mode string) error {
	switch strings.TrimSpace(mode) {
	case "", "auto", "xgo", "go":
		return nil
	default:
		return fmt.Errorf("invalid --go-cmd value %q: must be one of auto, xgo, go", mode)
	}
}
