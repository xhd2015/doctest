package leafcache

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// LeafRef identifies one leaf under an absolute (or abs-cleaned) tree root.
// Used by multi-tree prepare so identities cannot collide on shared relative paths.
type LeafRef struct {
	TreeRoot string
	LeafRel  string
}

// PassPlan is the in-memory prepare result: store keys keyed by leaf identity,
// plus warm skip identities when skip is enabled.
type PassPlan struct {
	// Keys maps FormatLeafIdentity(tree, rel) -> ComputeLeafKey hex.
	Keys map[string]string
	// Skip is sorted warm identities (GetPass hits) when skipEnabled; empty otherwise.
	Skip []string
}

// FormatLeafIdentity returns a stable, non-empty, tree-qualified identity token
// for skip lists and fail maps. Different absolute TreeRoots with the same
// relative leaf path produce different identities (bare relpaths collide).
//
// The in-memory form uses a NUL separator (cannot appear in paths). For process
// environment transport use FormatLeafIdentityEnv (tab separator).
func FormatLeafIdentity(treeRoot, leafRel string) string {
	root := strings.TrimSpace(treeRoot)
	if abs, err := filepath.Abs(root); err == nil {
		root = abs
	}
	root = filepath.ToSlash(filepath.Clean(root))
	rel := filepath.ToSlash(filepath.Clean(strings.TrimSpace(leafRel)))
	if rel == "." {
		rel = ""
	}
	// NUL separator: abs-root+rel cannot collide with nested path pairs
	// (e.g. /a + b/c vs /a/b + c).
	return root + "\x00" + rel
}

// FormatLeafIdentityEnv returns an env-safe form of FormatLeafIdentity for
// DOCTEST_LEAF_CACHE_SKIP_PATHS (NUL cannot appear in process environment values).
// Generated __wreg / RunAll match this encoding when abs TreeRoot is baked in.
func FormatLeafIdentityEnv(treeRoot, leafRel string) string {
	return strings.ReplaceAll(FormatLeafIdentity(treeRoot, leafRel), "\x00", "\t")
}

// IdentityToEnv converts an in-memory FormatLeafIdentity token to its env form.
func IdentityToEnv(id string) string {
	return strings.ReplaceAll(id, "\x00", "\t")
}

// PreparePassPlan computes store keys for each leaf and, when skipEnabled,
// collects warm GetPass identities into a sorted Skip list.
// When skipEnabled is false, Skip is empty/nil but Keys is still filled.
func PreparePassPlan(store *Store, leaves []LeafRef, goVersion string, skipEnabled bool) (PassPlan, error) {
	plan := PassPlan{
		Keys: make(map[string]string, len(leaves)),
	}
	if store == nil {
		return plan, fmt.Errorf("leafcache: nil store")
	}
	for _, leaf := range leaves {
		id := FormatLeafIdentity(leaf.TreeRoot, leaf.LeafRel)
		in, err := KeyForLeaf(leaf.TreeRoot, leaf.LeafRel, goVersion)
		if err != nil {
			return PassPlan{}, err
		}
		key, err := ComputeLeafKey(in)
		if err != nil {
			return PassPlan{}, err
		}
		plan.Keys[id] = key
		if !skipEnabled {
			continue
		}
		hit, err := store.GetPass(key)
		if err != nil {
			return PassPlan{}, err
		}
		if hit {
			plan.Skip = append(plan.Skip, id)
		}
	}
	sort.Strings(plan.Skip)
	return plan, nil
}

// RecordPasses writes PutPass for identities that should be cached as warm.
// When allPassed is true, every key is stored. Otherwise PutPass only for
// identities not marked in failed; when !allPassed && failed == nil, store nothing.
// Store I/O errors are ignored (best-effort).
func RecordPasses(store *Store, keys map[string]string, failed map[string]bool, allPassed bool) {
	if store == nil || len(keys) == 0 {
		return
	}
	for id, key := range keys {
		if !allPassed {
			if failed != nil && failed[id] {
				continue
			}
			// Partial fail without per-leaf map: avoid storing may-have-failed leaves.
			if failed == nil {
				continue
			}
		}
		_ = store.PutPass(key)
	}
}
