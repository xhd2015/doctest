// Package leafcache computes stable leaf source-hash keys and stores
// explicit pass markers on disk (P1: local content DAG + store).
package leafcache

// AlgoVersion is mixed into every leaf key so format bumps invalidate
// existing store entries.
const AlgoVersion = "v1"

// KeyInput identifies one leaf for ComputeLeafKey.
type KeyInput struct {
	// ModuleRoot is the absolute directory containing go.mod.
	ModuleRoot string
	// TreeRoot is the absolute directory containing root DOCTEST.md.
	TreeRoot string
	// LeafDir is the absolute leaf directory (ASSERT.md present).
	LeafDir string
	// GoVersion is the toolchain version string, e.g. "go1.25.0".
	GoVersion string
}
