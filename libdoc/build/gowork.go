package build

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/xhd2015/doctest/libdoc/core"
)

const (
	// HubDirName is the multi-module orchestration module under the toplevel gen root.
	HubDirName = "__hub"
	// hubModulePath is the module path of the hub go.mod.
	hubModulePath = "testcase/hub"
	// genModuleName is the default module path written by core.WriteGoMod.
	genModuleName = "testcase"
)

// uniqueWorkModulePath returns a stable unique module path for a gen root so
// multiple gen modules can be required/replaced from the hub (duplicate
// "testcase" module paths are illegal).
func uniqueWorkModulePath(genRoot string) string {
	sum := sha256.Sum256([]byte(filepath.Clean(genRoot)))
	return genModuleName + "/w" + hex.EncodeToString(sum[:8])
}

// pickToplevelGenRoot chooses the outermost gen root (path prefix of all others).
// If roots are siblings (no nesting), returns the longest common parent directory.
func pickToplevelGenRoot(genRoots []string) string {
	if len(genRoots) == 0 {
		return ""
	}
	cleaned := make([]string, len(genRoots))
	for i, r := range genRoots {
		cleaned[i] = filepath.Clean(r)
	}
	sort.Slice(cleaned, func(i, j int) bool {
		if len(cleaned[i]) != len(cleaned[j]) {
			return len(cleaned[i]) < len(cleaned[j])
		}
		return cleaned[i] < cleaned[j]
	})
	// Prefer a root that is a strict prefix of every other (nested modules).
	for _, cand := range cleaned {
		ok := true
		for _, o := range cleaned {
			if o == cand {
				continue
			}
			if o != cand && !strings.HasPrefix(o, cand+string(filepath.Separator)) {
				ok = false
				break
			}
		}
		if ok {
			return cand
		}
	}
	// Siblings: common parent of all.
	prefix := cleaned[0]
	for _, o := range cleaned[1:] {
		prefix = longestCommonDirPrefix(prefix, o)
	}
	if prefix == "" || prefix == "." || prefix == string(filepath.Separator) {
		return cleaned[0]
	}
	return prefix
}

func longestCommonDirPrefix(a, b string) string {
	a = filepath.Clean(a)
	b = filepath.Clean(b)
	as := strings.Split(a, string(filepath.Separator))
	bs := strings.Split(b, string(filepath.Separator))
	n := len(as)
	if len(bs) < n {
		n = len(bs)
	}
	var parts []string
	for i := 0; i < n; i++ {
		if as[i] != bs[i] {
			break
		}
		parts = append(parts, as[i])
	}
	if len(parts) == 0 {
		return ""
	}
	// Preserve leading slash on Unix absolute paths.
	out := filepath.Join(parts...)
	if strings.HasPrefix(a, string(filepath.Separator)) && !strings.HasPrefix(out, string(filepath.Separator)) {
		out = string(filepath.Separator) + out
	}
	return out
}

// writeFileIfChangedPlain writes only when content differs.
func writeFileIfChangedPlain(path string, data []byte) error {
	existing, err := os.ReadFile(path)
	if err == nil && string(existing) == string(data) {
		return nil
	}
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// ensureWorkModulePath rewrites genRoot's module path to uniquePath and rewrites
// import paths in .go files. Does not descend into subdirectories that contain
// their own go.mod (nested module gen roots).
func ensureWorkModulePath(genRoot, uniquePath string) error {
	goModPath := filepath.Join(genRoot, "go.mod")
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("read go.mod: %w", err)
	}
	curMod, ok := parseGoModModulePath(string(data))
	if !ok {
		return fmt.Errorf("go.mod: missing module line in %s", genRoot)
	}
	if curMod == uniquePath {
		return nil
	}
	if curMod != genModuleName && !strings.HasPrefix(curMod, genModuleName+"/w") {
		return fmt.Errorf("go.mod module %q not a gen module (want %s or %s/w…)", curMod, genModuleName, genModuleName)
	}

	newGoMod := strings.Replace(string(data), "module "+curMod, "module "+uniquePath, 1)
	if err := os.WriteFile(goModPath, []byte(newGoMod), 0o644); err != nil {
		return err
	}

	oldPrefix := curMod + "/"
	newPrefix := uniquePath + "/"
	return walkModuleGoFiles(genRoot, func(path string, info os.FileInfo) error {
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		s := string(b)
		ns := strings.ReplaceAll(s, `"`+oldPrefix, `"`+newPrefix)
		ns = strings.ReplaceAll(ns, `"`+curMod+`"`, `"`+uniquePath+`"`)
		if ns == s {
			return nil
		}
		return os.WriteFile(path, []byte(ns), info.Mode().Perm())
	})
}

// walkModuleGoFiles walks genRoot but skips subdirectories that contain go.mod
// (nested module boundaries), other than genRoot itself.
func walkModuleGoFiles(genRoot string, fn func(path string, info os.FileInfo) error) error {
	genRoot = filepath.Clean(genRoot)
	return filepath.Walk(genRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if path == genRoot {
				return nil
			}
			// Nested gen module (e.g. parent/sub with its own go.mod).
			if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
				return filepath.SkipDir
			}
			// Skip hub dir if present.
			if info.Name() == HubDirName {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		return fn(path, info)
	})
}

func parseGoModModulePath(content string) (string, bool) {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), true
		}
	}
	return "", false
}

// suiteImportForPrep returns the import path of the suite package (has RunAll)
// after ensureWorkModulePath.
func suiteImportForPrep(uniquePath, treeRel string, multiTree bool) string {
	if multiTree {
		return uniquePath + "/" + core.WorkspaceDirName + "/" + core.WorkspaceSuiteDirName
	}
	treeRel = filepath.ToSlash(filepath.Clean(treeRel))
	if treeRel == "" || treeRel == "." {
		return uniquePath + "/" + core.UnifiedSuiteDirName
	}
	return uniquePath + "/" + strings.TrimPrefix(treeRel, "./") + "/" + core.UnifiedSuiteDirName
}

// memberSuite describes one suite.RunAll import for the hub.
type memberSuite struct {
	Alias string // valid Go identifier
	Path  string // import path
	Name  string // subtest name (tree or module label)
}

// writeMultiModHub writes toplevel/__hub/go.mod + suite that calls each RunAll.
// Returns hub directory (cwd for go test). goCache is optional isolated GOCACHE
// for hub go mod tidy (cmd.Env only).
func writeMultiModHub(toplevel string, members []memberSuite, replaceByMod map[string]string, goCache string) (hubDir string, err error) {
	hubDir = filepath.Join(toplevel, HubDirName)
	if err := os.MkdirAll(filepath.Join(hubDir, "suite"), 0o755); err != nil {
		return "", err
	}

	// go.mod: require each member + replace members and copy member replaces
	// (e.g. github.com/xhd2015/doctest => …) so go test does not try to download.
	var mod strings.Builder
	mod.WriteString("module ")
	mod.WriteString(hubModulePath)
	mod.WriteString("\n\ngo 1.21\n\n")
	modPaths := make([]string, 0, len(replaceByMod))
	for m := range replaceByMod {
		modPaths = append(modPaths, m)
	}
	sort.Strings(modPaths)
	if len(modPaths) > 0 {
		mod.WriteString("require (\n")
		for _, m := range modPaths {
			mod.WriteString("\t")
			mod.WriteString(m)
			mod.WriteString(" v0.0.0\n")
		}
		mod.WriteString(")\n\n")
		for _, m := range modPaths {
			abs, aerr := filepath.Abs(replaceByMod[m])
			if aerr != nil {
				abs = replaceByMod[m]
			}
			mod.WriteString("replace ")
			mod.WriteString(m)
			mod.WriteString(" => ")
			mod.WriteString(filepath.ToSlash(abs))
			mod.WriteString("\n")
		}
		// Union of replace lines from member go.mod files (skip module testcase/w*).
		seenReplace := map[string]bool{}
		for _, m := range modPaths {
			extra, rerr := collectGoModReplaces(replaceByMod[m])
			if rerr != nil {
				continue
			}
			for _, line := range extra {
				// line is "path => target"
				key := strings.SplitN(line, " => ", 2)[0]
				if strings.HasPrefix(key, genModuleName) {
					continue
				}
				if seenReplace[key] {
					continue
				}
				seenReplace[key] = true
				mod.WriteString("replace ")
				mod.WriteString(line)
				mod.WriteString("\n")
			}
		}
	}
	if err := writeFileIfChangedPlain(filepath.Join(hubDir, "go.mod"), []byte(mod.String())); err != nil {
		return "", err
	}
	// Member gen roots (for merged xgo-style vendor go.mod overlay on hub tidy).
	memberGenRoots := make([]string, 0, len(replaceByMod))
	for _, g := range replaceByMod {
		if g != "" {
			memberGenRoots = append(memberGenRoots, g)
		}
	}
	// Populate go.sum / resolve graph for the hub module.
	if err := tidyHubGoMod(hubDir, goCache, memberGenRoots); err != nil {
		return "", fmt.Errorf("hub go mod tidy: %w", err)
	}

	// suite/runall.go — not needed; hub only has suite_test.go that calls members.
	// suite/suite_test.go
	var suite strings.Builder
	suite.WriteString("package suite\n\n")
	suite.WriteString("import (\n")
	suite.WriteString("\t\"testing\"\n\n")
	// sort members by name for stable output
	sorted := append([]memberSuite(nil), members...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Name < sorted[j].Name })
	for _, m := range sorted {
		suite.WriteString("\t")
		suite.WriteString(m.Alias)
		suite.WriteString(" \"")
		suite.WriteString(m.Path)
		suite.WriteString("\"\n")
	}
	suite.WriteString(")\n\n")
	suite.WriteString("func TestDoctestSuite(t *testing.T) {\n")
	for _, m := range sorted {
		suite.WriteString("\tt.Run(")
		suite.WriteString(fmt.Sprintf("%q", m.Name))
		suite.WriteString(", ")
		suite.WriteString(m.Alias)
		suite.WriteString(".RunAll)\n")
	}
	suite.WriteString("}\n")
	if err := core.WriteFormattedGo(filepath.Join(hubDir, "suite", "suite_test.go"), suite.String()); err != nil {
		return "", err
	}
	return hubDir, nil
}

// collectGoModReplaces returns "path => target" lines from genRoot/go.mod.
func collectGoModReplaces(genRoot string) ([]string, error) {
	data, err := os.ReadFile(filepath.Join(genRoot, "go.mod"))
	if err != nil {
		return nil, err
	}
	var out []string
	inBlock := false
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "replace (" {
			inBlock = true
			continue
		}
		if inBlock {
			if trimmed == ")" {
				inBlock = false
				continue
			}
			if strings.HasPrefix(trimmed, "replace ") {
				trimmed = strings.TrimPrefix(trimmed, "replace ")
			}
			if trimmed != "" {
				out = append(out, trimmed)
			}
			continue
		}
		if strings.HasPrefix(trimmed, "replace ") {
			out = append(out, strings.TrimPrefix(trimmed, "replace "))
		}
	}
	return out, nil
}

// tidyHubGoMod runs go mod tidy in hubDir. memberGenRoots supply per-member
// vendor-gomod-overlay.json maps, merged under hubDir so phantom vendor go.mod
// files are visible (replace targets project vendor without on-disk go.mod).
func tidyHubGoMod(hubDir, goCache string, memberGenRoots []string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = hubDir
	var envExtras []string
	if goCache != "" {
		envExtras = append(envExtras, "GOCACHE="+goCache)
	}
	overlayPath, err := core.MergeVendorGomodOverlays(hubDir, memberGenRoots)
	if err != nil {
		return fmt.Errorf("merge vendor-gomod overlays: %w", err)
	}
	if overlayPath != "" {
		envExtras = append(envExtras, core.AppendGOFLAGSOverlay(overlayPath)...)
	}
	if len(envExtras) > 0 {
		cmd.Env = core.ChildEnv(nil, envExtras...)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w\n%s", err, out)
	}
	return nil
}

// aliasForImportPath builds a short unique Go identifier from an import path.
func aliasForImportPath(importPath string, used map[string]bool) string {
	base := filepath.Base(importPath)
	if base == "" || base == "." || base == "suite" {
		// use last non-suite segment
		parts := strings.Split(importPath, "/")
		for i := len(parts) - 1; i >= 0; i-- {
			if parts[i] != "suite" && parts[i] != "__workspace" && parts[i] != "" {
				base = parts[i]
				break
			}
		}
	}
	base = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return '_'
	}, base)
	if base == "" || (base[0] >= '0' && base[0] <= '9') {
		base = "m_" + base
	}
	alias := base
	for n := 2; used[alias]; n++ {
		alias = fmt.Sprintf("%s%d", base, n)
	}
	used[alias] = true
	return alias
}
