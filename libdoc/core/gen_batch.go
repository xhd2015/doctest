package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// genRootWriteMu serializes generate+prune for a shared mapping-gen root so
// concurrent nested suite leaves do not corrupt each other's emit notes or
// prune each other's packages.
var genRootWriteMu sync.Map // abs gen root → *sync.Mutex

// LockGenRootWrites locks exclusive generate/prune for genRoot. Caller must
// unlock (typically via defer).
func LockGenRootWrites(genRoot string) (unlock func()) {
	key := absGenRoot(genRoot)
	v, _ := genRootWriteMu.LoadOrStore(key, &sync.Mutex{})
	mu := v.(*sync.Mutex)
	mu.Lock()
	return mu.Unlock
}

// Attach registers this batch as the active desired-path recorder for genRoot
// so WriteIfChanged / WriteGoMod can Note without threading Options.
// Caller should hold LockGenRootWrites for the generate duration.
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
// Only safe under LockGenRootWrites for that gen root.
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

// rootBookkeeping reports shared gen-root files that must never be deleted by
// tree-scoped orphan prune (multi-tree / nested suite leaves share go.mod).
func rootBookkeeping(rel string) bool {
	switch rel {
	case "go.mod", "go.sum", genManifestFile, "doctest.tidy-done":
		return true
	default:
		return strings.HasPrefix(rel, "doctest.")
	}
}

// nestedTreeExcludes returns prefixes of other treeRels that live under treeRel
// so a parent tree prune never deletes a nested tree's packages.
func nestedTreeExcludes(treeRel string, allTreeRels []string) []string {
	treeRel = filepath.ToSlash(filepath.Clean(treeRel))
	if treeRel == "" || treeRel == "." {
		return nil
	}
	prefix := treeRel + "/"
	var out []string
	for _, o := range allTreeRels {
		o = filepath.ToSlash(filepath.Clean(o))
		if o == "" || o == "." || o == treeRel {
			continue
		}
		if strings.HasPrefix(o+"/", prefix) || strings.HasPrefix(o, prefix) {
			out = append(out, o+"/")
		}
	}
	return out
}

// PruneTreeScopeToDesired deletes orphan files under one tree's ownership scope
// inside genRoot. treeRel is filepath.Rel(modRoot, doctestRoot) ("" or "." =
// packages at gen root). Shared mapping-gen holds many trees: we never delete
// outside this tree's prefix, never under nested sibling treeRels, and never
// delete root bookkeeping (go.mod, …).
//
// Safety: only delete files that appear in doctest.gen-manifest for this gen
// root and are not in desired. That way an incomplete emit set cannot remove
// live package sources that were never recorded (false orphans).
//
// excludeNested are other treeRel prefixes (with trailing slash) that must not
// be touched when pruning a parent tree (e.g. tests/ vs tests/boundary/).
//
// desired rels are slash form relative to genRoot. Empty desired → no-op.
func PruneTreeScopeToDesired(genRoot, treeRel string, desired map[string]struct{}, excludeNested []string) error {
	if genRoot == "" {
		return nil
	}
	genRoot = absGenRoot(genRoot)
	if len(desired) == 0 {
		return nil
	}

	keep := make(map[string]struct{}, len(desired)+4)
	for k := range desired {
		keep[filepath.ToSlash(k)] = struct{}{}
	}
	for _, name := range []string{"go.mod", "go.sum", genManifestFile, "doctest.tidy-done"} {
		keep[name] = struct{}{}
	}

	treeRel = filepath.ToSlash(filepath.Clean(treeRel))
	if treeRel == "" {
		treeRel = "."
	}

	var scopePrefix string
	var ownedTops map[string]struct{}
	if treeRel != "." {
		scopePrefix = treeRel + "/"
	} else {
		ownedTops = map[string]struct{}{}
		for rel := range keep {
			if rootBookkeeping(rel) {
				continue
			}
			top := rel
			if i := strings.IndexByte(rel, '/'); i >= 0 {
				top = rel[:i]
			}
			if top == WorkspaceDirName {
				continue
			}
			ownedTops[top] = struct{}{}
		}
	}

	underNested := func(rel string) bool {
		for _, ex := range excludeNested {
			if strings.HasPrefix(rel, ex) || rel+"/" == ex {
				return true
			}
			// exact nested treeRel without trailing slash
			if strings.TrimSuffix(ex, "/") == rel {
				return true
			}
		}
		return false
	}

	inScope := func(rel string) bool {
		if rootBookkeeping(rel) {
			return false
		}
		if underNested(rel) {
			return false
		}
		if scopePrefix != "" {
			return rel == treeRel || strings.HasPrefix(rel, scopePrefix)
		}
		if strings.HasPrefix(rel, WorkspaceDirName+"/") || rel == WorkspaceDirName {
			return false
		}
		top := rel
		if i := strings.IndexByte(rel, '/'); i >= 0 {
			top = rel[:i]
		}
		_, ok := ownedTops[top]
		return ok
	}

	genModMu.Lock()
	man, err := cachedGenManifestLocked(genRoot)
	if err != nil {
		genModMu.Unlock()
		return err
	}
	// Snapshot candidate rels from manifest only (safe orphan set).
	var candidates []string
	for rel := range man.hashes {
		if !inScope(rel) {
			continue
		}
		if _, ok := keep[rel]; ok {
			continue
		}
		candidates = append(candidates, rel)
	}
	for _, rel := range candidates {
		man.deleteHash(rel)
	}
	if man.dirty {
		if err := man.flush(genRoot); err != nil {
			genModMu.Unlock()
			return err
		}
	}
	genModMu.Unlock()

	for _, rel := range candidates {
		p := filepath.Join(genRoot, filepath.FromSlash(rel))
		if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("PruneTreeScopeToDesired: remove %s: %w", p, err)
		}
	}
	// Drop empty dirs under scope.
	for {
		removed := 0
		_ = filepath.Walk(genRoot, func(path string, info os.FileInfo, err error) error {
			if err != nil || path == genRoot || info == nil || !info.IsDir() {
				return nil
			}
			rel, rerr := filepath.Rel(genRoot, path)
			if rerr != nil {
				return nil
			}
			rel = filepath.ToSlash(rel)
			if scopePrefix != "" {
				if rel != treeRel && !strings.HasPrefix(rel, scopePrefix) {
					return nil
				}
			} else if !inScope(rel) && rel != treeRel {
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
	return nil
}

// PruneWorkspaceScope deletes orphans under __workspace/ only (multi-tree hub).
func PruneWorkspaceScope(genRoot string, desired map[string]struct{}) error {
	if genRoot == "" {
		return nil
	}
	wsDesired := make(map[string]struct{})
	prefix := WorkspaceDirName + "/"
	for rel := range desired {
		rel = filepath.ToSlash(rel)
		if rel == WorkspaceDirName || strings.HasPrefix(rel, prefix) {
			wsDesired[rel] = struct{}{}
		}
	}
	if len(wsDesired) == 0 {
		return nil
	}
	return PruneTreeScopeToDesired(genRoot, WorkspaceDirName, wsDesired, nil)
}

// PruneTree prunes one tree's package scope using this batch's desired set.
// allTreeRels lists every tree in the workspace batch so parent prune excludes
// nested sibling trees (tests/ must not delete tests/boundary/ packages).
func (b *GenBatch) PruneTree(genRoot, treeRel string, allTreeRels []string) error {
	if b == nil {
		return nil
	}
	return PruneTreeScopeToDesired(genRoot, treeRel, b.Desired(genRoot), nestedTreeExcludes(treeRel, allTreeRels))
}

// PruneWorkspace prunes __workspace using this batch's desired set.
func (b *GenBatch) PruneWorkspace(genRoot string) error {
	if b == nil {
		return nil
	}
	return PruneWorkspaceScope(genRoot, b.Desired(genRoot))
}
