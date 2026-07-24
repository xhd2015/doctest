package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// genSessionMarker is written under a gen root after a session wipe so multi-tree
// prepare into the same throwaway cache does not re-wipe mid-invocation.
const genSessionMarker = "doctest.gen-session"

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

// EnsureCleanGenRoot makes genRoot a fresh throwaway cache for this generate
// session when needed:
//
//   - When sessionID is non-empty: wipe if the gen root is missing the session
//     marker or the marker does not match (once per CLI session / multi-tree
//     batch). Subsequent trees sharing the same gen root skip the wipe.
//   - When forceA is true and sessionID is empty: wipe once using a fixed marker
//     so library callers with -a still get a clean root.
//   - When sessionID is empty and forceA is false: no wipe (library multi-tree
//     prepare tests share a gen dir without a session id).
//
// The gen root is always ensured to exist. The session marker is rewritten after
// a wipe so parallel/sequential prepare of other trees does not clear them.
func EnsureCleanGenRoot(genRoot, sessionID string, forceA bool) error {
	if genRoot == "" {
		return fmt.Errorf("EnsureCleanGenRoot: empty genRoot")
	}
	genRoot = filepath.Clean(genRoot)
	if abs, err := filepath.Abs(genRoot); err == nil {
		genRoot = abs
	}

	want := strings.TrimSpace(sessionID)
	if want == "" {
		if !forceA {
			if err := os.MkdirAll(genRoot, 0o755); err != nil {
				return err
			}
			return nil
		}
		// Library -a without session: wipe when marker absent/mismatched.
		want = "force-a"
	}

	markerPath := filepath.Join(genRoot, genSessionMarker)
	if b, err := os.ReadFile(markerPath); err == nil && string(b) == want {
		return nil
	}

	if err := WipeGenRoot(genRoot); err != nil {
		return err
	}
	if err := os.WriteFile(markerPath, []byte(want), 0o644); err != nil {
		return fmt.Errorf("EnsureCleanGenRoot: write session marker: %w", err)
	}
	return nil
}
