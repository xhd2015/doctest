package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Unified gen-root content-hash skip index (replaces doctest.gomod-fp).
const (
	genManifestFile    = "doctest.gen-manifest"
	genManifestVersion = 1
)

// genManifest is an in-memory path → content-hash map for one gen root.
// Paths are slash-separated relative to the gen root.
type genManifest struct {
	version int
	hashes  map[string]string
	dirty   bool
}

func contentSHA256(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func loadGenManifest(genRoot string) (*genManifest, error) {
	m := &genManifest{
		version: genManifestVersion,
		hashes:  make(map[string]string),
	}
	path := filepath.Join(genRoot, genManifestFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return m, nil
		}
		return nil, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		if fields[0] == "version" {
			var v int
			if _, err := fmt.Sscanf(fields[1], "%d", &v); err == nil {
				m.version = v
			}
			continue
		}
		// "path hash" — path may contain no spaces (slash-separated rel paths).
		rel := fields[0]
		hash := fields[1]
		m.hashes[rel] = hash
	}
	if m.version != genManifestVersion {
		// Format bump: treat as cold (force full re-hash misses).
		m.version = genManifestVersion
		m.hashes = make(map[string]string)
		// Not dirty until a write updates entries; next successful write flushes.
	}
	return m, nil
}

func (m *genManifest) setHash(rel, hash string) {
	rel = filepath.ToSlash(rel)
	if m.hashes[rel] == hash {
		return
	}
	m.hashes[rel] = hash
	m.dirty = true
}

func (m *genManifest) deleteHash(rel string) {
	rel = filepath.ToSlash(rel)
	if _, ok := m.hashes[rel]; !ok {
		return
	}
	delete(m.hashes, rel)
	m.dirty = true
}

func (m *genManifest) encode() []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "version %d\n", genManifestVersion)
	keys := make([]string, 0, len(m.hashes))
	for k := range m.hashes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(&b, "%s %s\n", k, m.hashes[k])
	}
	return []byte(b.String())
}

// flush writes the manifest only when the in-memory map changed (content-stable).
func (m *genManifest) flush(genRoot string) error {
	if !m.dirty {
		return nil
	}
	if err := os.MkdirAll(genRoot, 0755); err != nil {
		return err
	}
	path := filepath.Join(genRoot, genManifestFile)
	data := m.encode()
	if _, err := writeFileIfChanged(path, data, 0644); err != nil {
		return err
	}
	m.dirty = false
	return nil
}

// writeRelIfChanged implements content-hash skip for one relative path under genRoot.
// On hash hit (and target exists), skips reading and rewriting the target.
// On hash miss, writes when needed and records the hash. Caller must hold genModMu
// when concurrent writers share genRoot.
func (m *genManifest) writeRelIfChanged(genRoot, rel string, data []byte) (wrote bool, err error) {
	rel = filepath.ToSlash(rel)
	hash := contentSHA256(data)
	abs := filepath.Join(genRoot, filepath.FromSlash(rel))
	if m.hashes[rel] == hash {
		if _, err := os.Stat(abs); err == nil {
			return false, nil
		}
		// Manifest claims hash but file missing — rewrite.
	}
	if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
		return false, err
	}
	wrote, err = writeFileIfChanged(abs, data, 0644)
	if err != nil {
		return false, err
	}
	m.setHash(rel, hash)
	return wrote, nil
}

// WriteIfChanged is the public unified gen writer: hash desired content against
// doctest.gen-manifest; skip rewrite on hit; write + update entry on miss.
func WriteIfChanged(genRoot, rel string, data []byte) (wrote bool, err error) {
	if genRoot == "" || rel == "" {
		return false, fmt.Errorf("WriteIfChanged: genRoot and rel required")
	}
	genModMu.Lock()
	defer genModMu.Unlock()

	if err := os.MkdirAll(genRoot, 0755); err != nil {
		return false, err
	}
	man, err := loadGenManifest(genRoot)
	if err != nil {
		return false, err
	}
	wrote, err = man.writeRelIfChanged(genRoot, rel, data)
	if err != nil {
		return false, err
	}
	if err := man.flush(genRoot); err != nil {
		return wrote, err
	}
	return wrote, nil
}

// findGenRootWithManifest walks up from path's directory looking for
// doctest.gen-manifest. Returns gen root, slash-relative path, and true if found.
func findGenRootWithManifest(absPath string) (genRoot, rel string, ok bool) {
	absPath = filepath.Clean(absPath)
	dir := filepath.Dir(absPath)
	for {
		manPath := filepath.Join(dir, genManifestFile)
		if st, err := os.Stat(manPath); err == nil && !st.IsDir() {
			rel, err := filepath.Rel(dir, absPath)
			if err != nil || strings.HasPrefix(rel, "..") {
				return "", "", false
			}
			return dir, filepath.ToSlash(rel), true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", "", false
		}
		dir = parent
	}
}
