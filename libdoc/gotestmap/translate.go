// Package gotestmap defines the ideal mapping from doctest path args to go test
// commands under the current generation model (synthetic testcase modules).
//
// Perfect run is path-shaped, not hub:
//
//	(cd genOuter && go test ./tree/mid/...) && (cd genNested && go test ./...)
//
// Translate is intentionally not yet ideal (broken / stub) so tests stay RED
// until the product mapping is fixed.
package gotestmap

import (
	"fmt"
	"path"
	"sort"
	"strings"
)

// Cmd is one go test invocation: run with Dir as cwd (project-relative).
type Cmd struct {
	// Dir is the module root to cd into before go test (relative to project root).
	// "." = outer module.
	Dir string
	// Pattern is the go test package pattern, e.g. "./tree/mid/..." or "./...".
	Pattern string
}

// String formats as shell-ish: (cd Dir && go test Pattern) or go test Pattern when Dir is ".".
func (c Cmd) String() string {
	if c.Dir == "" || c.Dir == "." {
		return "go test " + c.Pattern
	}
	return "(cd " + path.Clean(c.Dir) + " && go test " + c.Pattern + ")"
}

// Layout describes fixture modules on disk (dirs that contain go.mod), relative to project root.
// ModuleRoots must include "." when the outer module exists at project root.
type Layout struct {
	ModuleRoots []string // e.g. ".", "tree/mid/nestedmod"
}

// Translate maps a doctest path arg to path-shaped go test commands.
//
// Rules (locked by tests):
//   - path/... → go test ./path/... in the module that owns path (outer)
//   - never expand to parent siblings (not go test ./tree/... when user said ./tree/mid/...)
//   - each nested go.mod under the path base → extra (cd nestedRoot && go test ./...)
//   - path without ... → go test ./path (single package pattern, no recursion)
func Translate(userArg string, layout Layout) ([]Cmd, error) {
	return IdealTranslate(userArg, layout)
}

// IdealTranslate is the correct mapping (for reference / future product).
// Tests assert against this; wire Translate to IdealTranslate to go green.
func IdealTranslate(userArg string, layout Layout) ([]Cmd, error) {
	arg := strings.TrimSpace(userArg)
	if arg == "" {
		return nil, fmt.Errorf("empty path")
	}
	if arg == "..." {
		return nil, fmt.Errorf("bare '...' not supported; use './...' or 'path/...'")
	}

	recursive := strings.HasSuffix(arg, "/...")
	base := arg
	if recursive {
		base = strings.TrimSuffix(arg, "/...")
	}
	base = strings.TrimPrefix(base, "./")
	if base == "" {
		base = "."
	}
	base = path.Clean(base)

	mods := normalizeModules(layout.ModuleRoots)

	// Outer module root that owns base (deepest module root that is a prefix of base).
	owner := owningModule(base, mods)
	// Pattern relative to owner module root.
	var pattern string
	if owner == "." {
		if base == "." {
			if recursive {
				pattern = "./..."
			} else {
				pattern = "."
			}
		} else if recursive {
			pattern = "./" + base + "/..."
		} else {
			pattern = "./" + base
		}
	} else {
		// base under nested module path owner
		rel := base
		if base == owner {
			rel = ""
		} else if strings.HasPrefix(base, owner+"/") {
			rel = strings.TrimPrefix(base, owner+"/")
		}
		if recursive {
			if rel == "" {
				pattern = "./..."
			} else {
				pattern = "./" + rel + "/..."
			}
		} else {
			if rel == "" {
				pattern = "."
			} else {
				pattern = "./" + rel
			}
		}
	}

	var out []Cmd
	out = append(out, Cmd{Dir: owner, Pattern: pattern})

	// Nested modules under path base only for recursive path/... (B/C).
	if recursive {
		prefix := base
		for _, m := range mods {
			if m == "." || m == owner {
				continue
			}
			if prefix == "." {
				// ./... → every other module root
				out = append(out, Cmd{Dir: m, Pattern: "./..."})
				continue
			}
			// Module root under selected base (including base itself as module root).
			if m == prefix || strings.HasPrefix(m, prefix+"/") {
				out = append(out, Cmd{Dir: m, Pattern: "./..."})
			}
		}
	}

	return dedupeCmds(out), nil
}

func normalizeModules(roots []string) []string {
	if len(roots) == 0 {
		return []string{"."}
	}
	seen := map[string]bool{}
	var out []string
	for _, r := range roots {
		r = path.Clean(strings.TrimPrefix(r, "./"))
		if r == "" {
			r = "."
		}
		if seen[r] {
			continue
		}
		seen[r] = true
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool {
		// deeper first for owningModule
		if strings.Count(out[i], "/") != strings.Count(out[j], "/") {
			return strings.Count(out[i], "/") > strings.Count(out[j], "/")
		}
		return out[i] < out[j]
	})
	return out
}

func owningModule(base string, mods []string) string {
	// mods sorted deeper first
	if base == "." {
		return "."
	}
	for _, m := range mods {
		if m == "." {
			continue
		}
		if base == m || strings.HasPrefix(base, m+"/") {
			return m
		}
	}
	return "."
}

func dedupeCmds(in []Cmd) []Cmd {
	type key struct{ d, p string }
	seen := map[key]bool{}
	var out []Cmd
	for _, c := range in {
		k := key{c.Dir, c.Pattern}
		if seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, c)
	}
	// stable: outer (Dir ".") first, then by Dir
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Dir == "." && out[j].Dir != "." {
			return true
		}
		if out[j].Dir == "." && out[i].Dir != "." {
			return false
		}
		if out[i].Dir != out[j].Dir {
			return out[i].Dir < out[j].Dir
		}
		return out[i].Pattern < out[j].Pattern
	})
	return out
}
