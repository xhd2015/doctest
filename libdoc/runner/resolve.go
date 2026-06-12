package runner

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/xhd2015/gitops/git"
)

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
		if path == moduleRoot {
			if hasFile(path, "DOCTEST.md") {
				dirs = append(dirs, path)
			}
			return nil
		}
		if hasFile(path, "go.mod") {
			nestedPath := readModulePath(path)
			if nestedPath == "" {
				return filepath.SkipDir
			}
			if !strings.HasPrefix(nestedPath, ancestorPath+"/") {
				return filepath.SkipDir
			}
		}
		if hasFile(path, "DOCTEST.md") {
			for _, existing := range dirs {
				if isAncestor(existing, path) {
					return nil
				}
			}
			dirs = append(dirs, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(dirs)
	return dirs, nil
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
		dir = parent
	}
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
		if ignored, err := git.CheckIgnore(absCwd, name); err == nil && ignored {
			continue
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

var errNoModuleFound = errors.New("no module found")

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
