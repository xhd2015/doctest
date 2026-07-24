package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// GenBatch tracks throwaway-gen lifecycle for one doctest test invocation
// (or library multi-tree prepare that shares the same *GenBatch pointer).
//
// It is intentionally unrelated to SessionID / DOCTEST_SESSION_ID.
//
// Multi-tree prepare must share one *GenBatch so:
//   - -a wipes each gen root at most once
//   - desired emit paths union across trees before orphan prune
type GenBatch struct {
	mu    sync.Mutex
	wiped map[string]struct{}            // abs gen root → wiped for -a this batch
	emit  map[string]map[string]struct{} // abs gen root → slash-rel desired paths
}

// NewGenBatch returns an empty batch.
func NewGenBatch() *GenBatch {
	return &GenBatch{
		wiped: make(map[string]struct{}),
		emit:  make(map[string]map[string]struct{}),
	}
}

// absGenRoot normalizes gen root for map keys.
func absGenRoot(genRoot string) string {
	genRoot = filepath.Clean(genRoot)
	if abs, err := filepath.Abs(genRoot); err == nil {
		return abs
	}
	return genRoot
}

// Attach registers this batch as the active desired-path recorder for genRoot
// so WriteIfChanged / WriteGoMod can Note without threading Options.
func (b *GenBatch) Attach(genRoot string) {
	if b == nil || genRoot == "" {
		return
	}
	activeBatches.Store(absGenRoot(genRoot), b)
}

// Detach removes the active registration for genRoot (if it still points at b).
func (b *GenBatch) Detach(genRoot string) {
	if b == nil || genRoot == "" {
		return
	}
	key := absGenRoot(genRoot)
	if cur, ok := activeBatches.Load(key); ok && cur == b {
		activeBatches.Delete(key)
	}
}

// activeBatches maps abs gen root → *GenBatch for NoteDesired during writes.
var activeBatches sync.Map

// NoteDesired records a path that this generate batch intends to keep under genRoot.
// rel is slash-separated relative to genRoot. No-op if no batch is attached.
func NoteDesired(genRoot, rel string) {
	if genRoot == "" || rel == "" {
		return
	}
	key := absGenRoot(genRoot)
	v, ok := activeBatches.Load(key)
	if !ok {
		return
	}
	b, _ := v.(*GenBatch)
	if b == nil {
		return
	}
	b.Note(key, rel)
}

// Note records a desired relative path under genRoot (already abs-clean preferred).
func (b *GenBatch) Note(genRoot, rel string) {
	if b == nil || genRoot == "" || rel == "" {
		return
	}
	genRoot = absGenRoot(genRoot)
	rel = filepath.ToSlash(filepath.Clean(rel))
	if rel == "." {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.emit == nil {
		b.emit = make(map[string]map[string]struct{})
	}
	m := b.emit[genRoot]
	if m == nil {
		m = make(map[string]struct{})
		b.emit[genRoot] = m
	}
	m[rel] = struct{}{}
}

// Desired returns a copy of desired rel paths for genRoot (may be nil/empty).
func (b *GenBatch) Desired(genRoot string) map[string]struct{} {
	if b == nil {
		return nil
	}
	genRoot = absGenRoot(genRoot)
	b.mu.Lock()
	defer b.mu.Unlock()
	src := b.emit[genRoot]
	if len(src) == 0 {
		return nil
	}
	out := make(map[string]struct{}, len(src))
	for k := range src {
		out[k] = struct{}{}
	}
	return out
}

// WipeOnce removes all contents of genRoot at most once per batch (CLI -a).
// Clears the emit set for that root so a fresh generate repopulates desired paths.
func (b *GenBatch) WipeOnce(genRoot string) error {
	if genRoot == "" {
		return fmt.Errorf("GenBatch.WipeOnce: empty genRoot")
	}
	genRoot = absGenRoot(genRoot)
	if b == nil {
		return WipeGenRoot(genRoot)
	}
	b.mu.Lock()
	if _, ok := b.wiped[genRoot]; ok {
		b.mu.Unlock()
		return os.MkdirAll(genRoot, 0o755)
	}
	if b.wiped == nil {
		b.wiped = make(map[string]struct{})
	}
	b.wiped[genRoot] = struct{}{}
	delete(b.emit, genRoot)
	b.mu.Unlock()

	if err := WipeGenRoot(genRoot); err != nil {
		return err
	}
	// Re-attach after wipe (same batch).
	b.Attach(genRoot)
	return nil
}

// WipeGenRoot removes all contents of genRoot (throwaway generate cache),
// recreates the directory, and drops the in-memory gen-manifest cache.
func WipeGenRoot(genRoot string) error {
	if genRoot == "" {
		return fmt.Errorf("WipeGenRoot: empty genRoot")
	}
	genRoot = filepath.Clean(genRoot)
	InvalidateGenManifestCache(genRoot)
	if err := os.RemoveAll(genRoot); err != nil {
		return fmt.Errorf("WipeGenRoot: remove %s: %w", genRoot, err)
	}
	if err := os.MkdirAll(genRoot, 0o755); err != nil {
		return fmt.Errorf("WipeGenRoot: mkdir %s: %w", genRoot, err)
	}
	return nil
}

// PruneGenRootToDesired deletes files under genRoot that are not in desired
// (orphan reconcile). Always retains doctest.gen-manifest as a bookkeeping file
// when present after rewrite. desired rels use slash form.
//
// If desired is empty, prune is a no-op (avoids accidental full delete when
// no generate notes were recorded).
func PruneGenRootToDesired(genRoot string, desired map[string]struct{}) error {
	if genRoot == "" {
		return nil
	}
	genRoot = absGenRoot(genRoot)
	if len(desired) == 0 {
		return nil
	}
	// Bookkeeping always kept if we prune.
	keep := make(map[string]struct{}, len(desired)+2)
	for k := range desired {
		keep[filepath.ToSlash(k)] = struct{}{}
	}
	keep[genManifestFile] = struct{}{}

	var toDelete []string
	err := filepath.Walk(genRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if path == genRoot {
			return nil
		}
		rel, rerr := filepath.Rel(genRoot, path)
		if rerr != nil {
			return rerr
		}
		rel = filepath.ToSlash(rel)
		if info.IsDir() {
			return nil // delete files first; prune empty dirs after
		}
		if _, ok := keep[rel]; ok {
			return nil
		}
		toDelete = append(toDelete, path)
		return nil
	})
	if err != nil {
		return fmt.Errorf("PruneGenRootToDesired: walk %s: %w", genRoot, err)
	}
	for _, p := range toDelete {
		if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("PruneGenRootToDesired: remove %s: %w", p, err)
		}
	}
	// Remove empty directories (bottom-up via Walk is messy; second walk).
	for {
		removed := 0
		_ = filepath.Walk(genRoot, func(path string, info os.FileInfo, err error) error {
			if err != nil || path == genRoot || info == nil || !info.IsDir() {
				return nil
			}
			entries, rerr := os.ReadDir(path)
			if rerr != nil || len(entries) > 0 {
				return nil
			}
			if err := os.Remove(path); err == nil {
				removed++
			}
			return nil
		})
		if removed == 0 {
			break
		}
	}

	// Drop stale manifest hashes for deleted rels.
	genModMu.Lock()
	defer genModMu.Unlock()
	man, err := cachedGenManifest(genRoot)
	if err != nil {
		return err
	}
	for rel := range man.hashes {
		if _, ok := keep[rel]; !ok {
			man.deleteHash(rel)
		}
	}
	if man.dirty {
		if err := man.flush(genRoot); err != nil {
			return err
		}
	}
	return nil
}

// PruneAttached prunes genRoot using this batch's desired set (copy under lock).
func (b *GenBatch) PruneAttached(genRoot string) error {
	if b == nil {
		return nil
	}
	return PruneGenRootToDesired(genRoot, b.Desired(genRoot))
}
