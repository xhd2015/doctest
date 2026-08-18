package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xhd2015/xgo/support/coverage"
)

// CoverFileKeysFromExposeMaterialized returns import-path cover file keys for
// session-generated expose.go files listed under genRoot/ExposeMaterializedList.
// Empty when the list is missing or empty. Keys are suitable for
// strings.HasPrefix against coverage.CovLine.Prefix (file path before range).
func CoverFileKeysFromExposeMaterialized(genRoot string) ([]string, error) {
	if genRoot == "" {
		return nil, nil
	}
	listPath := filepath.Join(genRoot, ExposeMaterializedList)
	data, err := os.ReadFile(listPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var keys []string
	seen := map[string]struct{}{}
	for _, line := range strings.Split(string(data), "\n") {
		path := strings.TrimSpace(line)
		if path == "" {
			continue
		}
		key, kerr := coverFileKeyFromExposeAbs(path)
		if kerr != nil {
			return nil, kerr
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		keys = append(keys, key)
	}
	return keys, nil
}

// StripExposeFromCoverProfile rewrites coverProfile in place, dropping lines
// whose Prefix matches any cover file key derived from ExposeMaterializedList
// under genRoots. No-op when coverProfile is empty/missing or no keys exist.
// Preserves the profile's coverage mode. Does not invent keys from reserved
// path name patterns alone.
func StripExposeFromCoverProfile(coverProfile string, genRoots []string) (removed int, err error) {
	if coverProfile == "" || len(genRoots) == 0 {
		return 0, nil
	}
	if _, statErr := os.Stat(coverProfile); statErr != nil {
		if os.IsNotExist(statErr) {
			return 0, nil
		}
		return 0, statErr
	}
	var keys []string
	seen := map[string]struct{}{}
	for _, root := range genRoots {
		part, kerr := CoverFileKeysFromExposeMaterialized(root)
		if kerr != nil {
			return 0, kerr
		}
		for _, k := range part {
			if _, ok := seen[k]; ok {
				continue
			}
			seen[k] = struct{}{}
			keys = append(keys, k)
		}
	}
	if len(keys) == 0 {
		return 0, nil
	}

	content, err := os.ReadFile(coverProfile)
	if err != nil {
		return 0, err
	}
	mode, lines := coverage.Parse(string(content))
	if mode == "" {
		mode = "set"
	}
	before := len(lines)
	lines = coverage.Filter(lines, func(line *coverage.CovLine) bool {
		if line == nil {
			return false
		}
		for _, k := range keys {
			if strings.HasPrefix(line.Prefix, k) {
				return false
			}
		}
		return true
	})
	removed = before - len(lines)
	if removed == 0 {
		return 0, nil
	}
	out := coverage.Format(mode, lines)
	if !strings.HasSuffix(out, "\n") {
		out += "\n"
	}
	if err := os.WriteFile(coverProfile, []byte(out), 0644); err != nil {
		return removed, err
	}
	return removed, nil
}

// coverFileKeyFromExposeAbs maps an absolute product expose.go path to the
// import-path file key used in go coverprofiles.
func coverFileKeyFromExposeAbs(abs string) (string, error) {
	if !isExposeMaterializedPath(abs) {
		return "", fmt.Errorf("cover file key: not an expose path: %s", abs)
	}
	dir := filepath.Dir(abs)
	for filepath.Base(dir) != DoctestInternalExposeDir {
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("cover file key: missing %s in %s", DoctestInternalExposeDir, abs)
		}
		dir = parent
	}
	productDir := filepath.Dir(dir)
	rel, err := filepath.Rel(productDir, abs)
	if err != nil {
		return "", fmt.Errorf("cover file key: rel %s under %s: %w", abs, productDir, err)
	}
	modPath, err := modulePathForDir(productDir)
	if err != nil {
		return "", err
	}
	return modPath + "/" + filepath.ToSlash(rel), nil
}

// modulePathForDir reads the module path from go.mod at dir or an ancestor.
func modulePathForDir(dir string) (string, error) {
	cur := filepath.Clean(dir)
	for {
		data, err := os.ReadFile(filepath.Join(cur, "go.mod"))
		if err == nil {
			if mod, ok := parseModulePathLine(string(data)); ok {
				return mod, nil
			}
			return "", fmt.Errorf("cover file key: no module line in %s/go.mod", cur)
		}
		if !os.IsNotExist(err) {
			return "", err
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			return "", fmt.Errorf("cover file key: no go.mod above %s", dir)
		}
		cur = parent
	}
}

func parseModulePathLine(content string) (string, bool) {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), true
		}
	}
	return "", false
}
