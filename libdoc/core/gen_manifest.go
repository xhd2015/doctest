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
	version      int
	hashes       map[string]string
	dirty        bool
	pendingFlush int // dirty updates since last flush
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
			// Still desired this batch even when rewrite skipped.
			NoteDesired(genRoot, rel)
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
	NoteDesired(genRoot, rel)
	return wrote, nil
}

// In-memory gen-manifest cache (gen root abs path → man). Avoids re-reading
// doctest.gen-manifest on every WriteIfChanged. Protected by genModMu.
var manByRoot = map[string]*genManifest{}

// WriteIfChanged is the public unified gen writer: hash desired content against
// doctest.gen-manifest; skip rewrite on hit; write + update entry on miss.
// Flushes the manifest at most every manifestFlushEvery updates (or when the
// entry is new) so parallel generate does not rewrite the whole index file
// once per Go package.
func WriteIfChanged(genRoot, rel string, data []byte) (wrote bool, err error) {
	if genRoot == "" || rel == "" {
		return false, fmt.Errorf("WriteIfChanged: genRoot and rel required")
	}
	genRoot = filepath.Clean(genRoot)
	genModMu.Lock()
	defer genModMu.Unlock()

	if err := os.MkdirAll(genRoot, 0755); err != nil {
		return false, err
	}
	man, err := cachedGenManifest(genRoot)
	if err != nil {
		return false, err
	}
	wrote, err = man.writeRelIfChanged(genRoot, rel, data)
	if err != nil {
		return false, err
	}
	// Always flush when dirty so callers (and nested selftests) reading
	// doctest.gen-manifest from disk see updates immediately after WriteIfChanged.
	// In-memory manByRoot still avoids re-reading the file on every call.
	if man.dirty {
		if err := man.flush(genRoot); err != nil {
			return wrote, err
		}
		man.pendingFlush = 0
	}
	return wrote, nil
}

// FlushGenManifest writes any dirty in-memory manifest for genRoot. Call after
// a tree's generate batch so hashes are durable before go test.
func FlushGenManifest(genRoot string) error {
	if genRoot == "" {
		return nil
	}
	genRoot = filepath.Clean(genRoot)
	genModMu.Lock()
	defer genModMu.Unlock()
	man := manByRoot[genRoot]
	if man == nil || !man.dirty {
		return nil
	}
	if err := man.flush(genRoot); err != nil {
		return err
	}
	man.pendingFlush = 0
	return nil
}

// InvalidateGenManifestCache drops the cached manifest (e.g. after cold wipe).
func InvalidateGenManifestCache(genRoot string) {
	if genRoot == "" {
		return
	}
	genRoot = filepath.Clean(genRoot)
	genModMu.Lock()
	delete(manByRoot, genRoot)
	genModMu.Unlock()
}

func cachedGenManifest(genRoot string) (*genManifest, error) {
	if man, ok := manByRoot[genRoot]; ok {
		return man, nil
	}
	man, err := loadGenManifest(genRoot)
	if err != nil {
		return nil, err
	}
	manByRoot[genRoot] = man
	return man, nil
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
