package path_resolve

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/xhd2015/gitops/git"
)

var ErrNoTestsFound = errors.New("no tests")

var errNoModuleFound = errors.New("no module found")

func IsDotDotDotPattern(arg string) bool {
	return strings.HasSuffix(arg, "/...")
}

func ExtractBasePath(arg string) string {
	base := strings.TrimSuffix(arg, "/...")
	base = strings.TrimPrefix(base, "./")
	if base == "" {
		return "."
	}
	return base
}

// defaultTreeWorkers caps concurrent light doctest trees for ./... runs.
const defaultTreeWorkers = 4

// RunForDirs discovers doctest trees under basePath and invokes fn for each.
// Independent trees run concurrently (bounded) to cut wall time for module-wide
// `./...` self-tests without thrashing the machine.
func RunForDirs(basePath string, fn func(dir string) error) error {
	n := defaultTreeWorkers
	if max := runtime.GOMAXPROCS(0); max > 0 && max < n {
		n = max
	}
	return RunForDirsLimit(basePath, n, fn)
}

// RunForDirsLimit is like RunForDirs but caps concurrent tree workers.
// workers <= 1 runs trees serially (stable for debugging).
func RunForDirsLimit(basePath string, workers int, fn func(dir string) error) error {
	dirs, err := FindDotDotDotDirs(basePath)
	if err != nil {
		return err
	}
	if len(dirs) == 0 {
		return ErrNoTestsFound
	}
	// Two-phase scheduling: light trees first (high fan-out), then heavy
	// nested-CLI trees (serialized). Mixing them inflates the main tests/
	// tree from ~2m to 4m+ via CPU thrash.
	var light, heavy []string
	for _, d := range dirs {
		if isHeavySelftestTree(d) {
			heavy = append(heavy, d)
		} else {
			light = append(light, d)
		}
	}
	sort.SliceStable(light, func(i, j int) bool {
		return estimateTreeWeight(light[i]) > estimateTreeWeight(light[j])
	})
	sort.SliceStable(heavy, func(i, j int) bool {
		return estimateTreeWeight(heavy[i]) > estimateTreeWeight(heavy[j])
	})

	var errs []string
	appendErr := func(dir string, err error) {
		if err == nil || errors.Is(err, ErrNoTestsFound) {
			return
		}
		errs = append(errs, dir+": "+err.Error())
	}

	// Phase 1: light trees in parallel.
	if err := runDirsParallel(light, workers, fn, appendErr); err != nil {
		return err
	}
	// Phase 2: heavy trees one at a time (full CPU for nested self-tests).
	for _, dir := range heavy {
		appendErr(dir, fn(dir))
	}
	if len(errs) > 0 {
		sort.Strings(errs)
		return errors.New("test failures:\n" + strings.Join(errs, "\n"))
	}
	return nil
}

func runDirsParallel(dirs []string, workers int, fn func(string) error, onErr func(string, error)) error {
	if len(dirs) == 0 {
		return nil
	}
	if workers <= 1 || len(dirs) == 1 {
		for _, dir := range dirs {
			onErr(dir, fn(dir))
		}
		return nil
	}
	if workers > len(dirs) {
		workers = len(dirs)
	}
	jobs := make(chan string, len(dirs))
	for _, d := range dirs {
		jobs <- d
	}
	close(jobs)
	var (
		mu sync.Mutex
		wg sync.WaitGroup
	)
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for dir := range jobs {
				err := fn(dir)
				if err == nil || errors.Is(err, ErrNoTestsFound) {
					continue
				}
				mu.Lock()
				onErr(dir, err)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return nil
}

// isHeavySelftestTree reports trees under the module's integration suite
// (…/doctest/tests/…), which shell out to the doctest binary extensively.
// Pure assert/libdoc/session trees are light and may run freely in parallel.
func isHeavySelftestTree(dir string) bool {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return false
	}
	// Match …/doctest/tests and …/doctest/tests/…
	sep := string(filepath.Separator)
	marker := sep + "doctest" + sep + "tests"
	if strings.HasSuffix(abs, marker) || strings.Contains(abs, marker+sep) {
		return true
	}
	return false
}

// estimateTreeWeight approximates tree cost by counting ASSERT.md leaves (cheap walk).
func estimateTreeWeight(dir string) int {
	n := 0
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if d.Name() == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() == "ASSERT.md" {
			n++
		}
		return nil
	})
	return n
}

func FindDotDotDotDirs(basePath string) ([]string, error) {
	if basePath == "" || basePath == "." {
		dirs, err := FindDOCTestDirsWithBase(".", ".")
		if err == nil {
			if len(dirs) == 0 {
				return nil, errors.New("no tests")
			}
			return dirs, nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		return FindDOCTestDirsFromSubdirs(".")
	}
	dirs, err := FindDOCTestDirsWithBase(basePath, basePath)
	absBase, absErr := filepath.Abs(basePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if absErr == nil {
				if _, ok := ResolveRoot(absBase); ok {
					return []string{absBase}, nil
				}
			}
			// Fall through to WalkDir below to find DOCTEST.md
			// directories outside of Go module hierarchy
			dirs = nil
		} else {
			return nil, err
		}
	}
	if absErr == nil {
		hasAbsBase := false
		for _, d := range dirs {
			if d == absBase {
				hasAbsBase = true
				break
			}
		}
		if !hasAbsBase {
			if _, ok := ResolveRoot(absBase); ok {
				dirs = append(dirs, absBase)
			}
		}

		filepath.WalkDir(absBase, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return nil
			}
			if !d.IsDir() || path == absBase {
				return nil
			}
			if d.Name() == "testdata" {
				return filepath.SkipDir
			}
			if hasFile(path, "DOCTEST.md") {
				if !containsDir(dirs, path) {
					dirs = append(dirs, path)
				}
				return filepath.SkipDir
			}
			return nil
		})
	}

	if len(dirs) == 0 {
		if absErr == nil {
			if _, ok := ResolveRoot(absBase); ok {
				return []string{absBase}, nil
			}
		}
		return nil, errors.New("no tests")
	}
	return dirs, nil
}

func ResolveRoot(dir string) (root string, ok bool) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", false
	}
	if _, err := os.Stat(absDir); err != nil {
		return "", false
	}

	var parents []string
	parents = append(parents, absDir)

	if !hasFile(absDir, "go.mod") {
		cur := absDir
		for {
			parent := filepath.Dir(cur)
			if parent == cur {
				break
			}
			if hasFile(parent, "go.mod") {
				parents = append(parents, parent)
				break
			}
			if filepath.Base(parent) == "testdata" {
				parents = append(parents, parent)
				break
			}
			parents = append(parents, parent)
			cur = parent
		}
	}

	for _, p := range parents {
		if hasFile(p, "DOCTEST.md") {
			return p, true
		}
	}

	for i := len(parents) - 1; i >= 0; i-- {
		if hasFile(parents[i], "SETUP.md") {
			return parents[i], true
		}
	}

	return "", false
}

func hasFile(dir, name string) bool {
	_, err := os.Stat(filepath.Join(dir, name))
	return err == nil
}

func FindDOCTestDirs(cwd string) ([]string, error) {
	moduleRoot, ancestorPath, err := findModuleRoot(cwd)
	if err != nil {
		return nil, err
	}

	var dirs []string
	err = filepath.WalkDir(moduleRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !d.IsDir() {
			return nil
		}
		if d.Name() == "testdata" {
			return filepath.SkipDir
		}
		if path == moduleRoot {
			if hasFile(path, "DOCTEST.md") {
				dirs = append(dirs, path)
			}
			return nil
		}
		if hasFile(path, ".git") {
			nestedPath := readModulePath(path)
			if nestedPath != "" {
				ancestorGit := gitRoot(moduleRoot)
				nestedGit := gitRoot(path)
				if ancestorGit != nestedGit {
					reason := gitSkipReason(ancestorGit, nestedGit)
					warnSkippingModule(nestedPath, path, ancestorPath, reason)
				}
			}
			return filepath.SkipDir
		}
		if hasFile(path, "go.mod") {
			nestedPath := readModulePath(path)
			if nestedPath == "" {
				return filepath.SkipDir
			}
			if !strings.HasPrefix(nestedPath, ancestorPath+"/") {
				ancestorGit := gitRoot(moduleRoot)
				nestedGit := gitRoot(path)
				if ancestorGit == nestedGit && shouldDiscoverNonChildModule(ancestorPath, nestedPath) {
					nestedDirs, nestedErr := FindDOCTestDirs(path)
					if nestedErr != nil {
						return nestedErr
					}
					dirs = append(dirs, nestedDirs...)
					return filepath.SkipDir
				}
				if ancestorGit != nestedGit {
					reason := gitSkipReason(ancestorGit, nestedGit)
					warnSkippingModule(nestedPath, path, ancestorPath, reason)
				}
				return filepath.SkipDir
			}
		}
		if hasFile(path, "DOCTEST.md") {
			if !containsDir(dirs, path) {
				dirs = append(dirs, path)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(dirs)
	return dirs, nil
}

func gitRoot(dir string) string {
	if inside, _ := git.IsInsideGit(dir); inside {
		root, err := git.ShowToplevel(dir)
		if err == nil {
			return root
		}
	}
	return ""
}

func shouldDiscoverNonChildModule(ancestorPath, nestedPath string) bool {
	family := moduleFamilyPrefix(ancestorPath)
	if strings.HasSuffix(family, "/") {
		return strings.HasPrefix(nestedPath, family)
	}
	return strings.HasPrefix(nestedPath, family)
}

func moduleFamilyPrefix(modulePath string) string {
	if i := strings.LastIndex(modulePath, "/"); i >= 0 {
		return modulePath[:i+1]
	}
	return modulePath
}

func gitSkipReason(ancestorGit, nestedGit string) string {
	if ancestorGit != "" && nestedGit != "" {
		return "different git repository"
	}
	return "git repository mismatch"
}

func warnSkippingModule(nestedPath, nestedDir, ancestorPath, reason string) {
	if ancestorPath != "" && nestedPath != "" && !strings.HasPrefix(nestedPath, ancestorPath+"/") {
		fmt.Fprintf(os.Stderr, "warning: skipping module %s at %s: not a child of %s (%s)\n",
			nestedPath, nestedDir, ancestorPath, reason)
		return
	}
	if nestedPath != "" {
		fmt.Fprintf(os.Stderr, "warning: skipping module %s at %s: %s\n",
			nestedPath, nestedDir, reason)
	}
}

func readModulePath(dir string) string {
	data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return ""
}

func findModuleRoot(cwd string) (dir string, modulePath string, err error) {
	dir, err = filepath.Abs(cwd)
	if err != nil {
		return "", "", err
	}
	gitRoot, gitRootErr := "", errors.New("not inside git")
	if inside, _ := git.IsInsideGit(dir); inside {
		gitRoot, gitRootErr = git.ShowToplevel(dir)
	}
	for {
		modFile := filepath.Join(dir, "go.mod")
		data, readErr := os.ReadFile(modFile)
		if readErr == nil {
			for _, line := range strings.Split(string(data), "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "module ") {
					modulePath = strings.TrimSpace(strings.TrimPrefix(line, "module "))
					return dir, modulePath, nil
				}
			}
			return dir, "", nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", "", os.ErrNotExist
		}
		if gitRootErr == nil && !isAncestor(gitRoot, parent) {
			return "", "", os.ErrNotExist
		}
		dir = parent
	}
}

func containsDir(dirs []string, dir string) bool {
	for _, existing := range dirs {
		if existing == dir {
			return true
		}
	}
	return false
}

func isAncestor(ancestor, child string) bool {
	rel, err := filepath.Rel(ancestor, child)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, "..")
}

func FindDOCTestDirsFromSubdirs(cwd string) ([]string, error) {
	absCwd, err := filepath.Abs(cwd)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(absCwd)
	if err != nil {
		return nil, err
	}

	var allDirs []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if inside, _ := git.IsInsideGit(absCwd); inside {
			if ignored, err := git.CheckIgnore(absCwd, name); err == nil && ignored {
				continue
			}
		}
		subdir := filepath.Join(absCwd, name)
		dirs, err := findDOCTestDirsInTree(subdir)
		if err != nil {
			if errors.Is(err, errNoModuleFound) {
				continue
			}
			return nil, err
		}
		allDirs = append(allDirs, dirs...)
	}

	sort.Strings(allDirs)
	return allDirs, nil
}

func findDOCTestDirsInTree(subdir string) ([]string, error) {
	var dirs []string
	err := filepath.WalkDir(subdir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !d.IsDir() {
			return nil
		}
		if hasFile(path, "go.mod") {
			nested, nestedErr := FindDOCTestDirs(path)
			if nestedErr != nil {
				return nestedErr
			}
			dirs = append(dirs, nested...)
			return filepath.SkipDir
		}
		if hasFile(path, "DOCTEST.md") {
			dirs = append(dirs, path)
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(dirs) == 0 {
		return nil, errNoModuleFound
	}
	return dirs, nil
}

func FindDOCTestDirsWithBase(cwd, basePath string) ([]string, error) {
	dirs, err := FindDOCTestDirs(cwd)
	if err != nil {
		return nil, err
	}
	if basePath == "" {
		return dirs, nil
	}
	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}
	var filtered []string
	for _, dir := range dirs {
		if isAncestor(absBase, dir) || dir == absBase {
			filtered = append(filtered, dir)
		}
	}
	return filtered, nil
}
