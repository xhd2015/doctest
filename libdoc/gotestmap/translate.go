// Package gotestmap is the single place that decides go test Dir+Pattern plans.
//
// Discovery/prepare/gen stay elsewhere. At the go-test call site, Plan() yields
// one or more Cmds.
//
// Production go-test path (Phase 1, wired via Plan into the workspace runner):
//
//   - ModeWorkspaceSuite: (cd gen && go test ./__workspace/suite)
//   - ModeHubSuite:       (cd hub && go test ./suite)
//
// Phase 2 (plan/contract only — not wired into the workspace runner):
//
//   - ModePathShaped / NeedsPathShaped / TranslatePath multi-cmd: mid-path and
//     cross-go.mod fixture cmds. Example for mid + nested module:
//     go test ./tree/mid/...
//     (cd tree/mid/nestedmod && go test ./...)
//
// finishWorkspaceGoTestCmds requires len(cmds)==1; multi-cmd path-shaped
// execution is deferred to a Phase 2 path-shaped executor.
package gotestmap

import (
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// Cmd is one go test invocation: run with Dir as cwd.
type Cmd struct {
	// Dir is the cwd for go test (project-relative, or absolute gen/hub path).
	// "." = outer project module (path-shaped fixture mapping).
	Dir string
	// Pattern is the go test package pattern, e.g. "./__workspace/suite", "./suite", "./tree/mid/...".
	Pattern string
}

// String formats as shell-ish: (cd Dir && go test Pattern) or go test Pattern when Dir is ".".
func (c Cmd) String() string {
	if c.Dir == "" || c.Dir == "." {
		return "go test " + c.Pattern
	}
	return "(cd " + c.Dir + " && go test " + c.Pattern + ")"
}

// Layout describes fixture modules on disk (dirs that contain go.mod), relative to project root.
// ModuleRoots must include "." when the outer module exists at project root.
type Layout struct {
	ModuleRoots []string // e.g. ".", "tree/mid/nestedmod"
}

// Mode selects which go-test plan family Plan returns.
type Mode int

const (
	// ModeWorkspaceSuite is today's single-gen fan-in: go test ./__workspace/suite (or equivalent).
	// Production go-test path uses this mode (or ModeHubSuite) via Plan.
	ModeWorkspaceSuite Mode = iota
	// ModeHubSuite is today's multi-mod hub: go test ./suite under __hub.
	// Production go-test path uses this mode (or ModeWorkspaceSuite) via Plan.
	ModeHubSuite
	// ModePathShaped is Phase 2 (plan/contract only): mid-leaf and/or cross-go.mod
	// fixture path patterns from TranslatePath. May return multiple Cmds.
	// Not wired into the workspace runner; do not pass multi-cmd results to
	// finishWorkspaceGoTestCmds (which requires len==1).
	ModePathShaped
)

// PlanInput is everything needed to produce go test cmds at the run site.
type PlanInput struct {
	Mode Mode

	// --- ModeWorkspaceSuite / ModeHubSuite (absolute run dirs) ---
	// RunDir is the process cwd for go test (gen root or hub dir). Absolute preferred.
	RunDir string
	// SuitePattern is the package arg, e.g. "./__workspace/suite" or "./suite".
	SuitePattern string

	// --- ModePathShaped (fixture-relative) ---
	UserArg string
	Layout  Layout
}

// Plan returns go test commands for the run site.
// Default suite/hub modes preserve today's single-command shape (production path).
// ModePathShaped uses TranslatePath (mid / nested modules) — Phase 2 plan/contract
// only; production runner does not execute multi-cmd path-shaped plans yet.
func Plan(in PlanInput) ([]Cmd, error) {
	switch in.Mode {
	case ModeWorkspaceSuite, ModeHubSuite:
		if strings.TrimSpace(in.RunDir) == "" {
			return nil, fmt.Errorf("gotestmap: RunDir required for suite/hub mode")
		}
		pat := strings.TrimSpace(in.SuitePattern)
		if pat == "" {
			if in.Mode == ModeHubSuite {
				pat = "./suite"
			} else {
				pat = "./__workspace/suite"
			}
		}
		return []Cmd{{Dir: in.RunDir, Pattern: pat}}, nil
	case ModePathShaped:
		return TranslatePath(in.UserArg, in.Layout)
	default:
		return nil, fmt.Errorf("gotestmap: unknown mode %d", in.Mode)
	}
}

// Translate is an alias for TranslatePath (fixture path → path-shaped cmds).
func Translate(userArg string, layout Layout) ([]Cmd, error) {
	return TranslatePath(userArg, layout)
}

// IdealTranslate is an alias for TranslatePath (historical test name).
func IdealTranslate(userArg string, layout Layout) ([]Cmd, error) {
	return TranslatePath(userArg, layout)
}

// TranslatePath maps a doctest path arg to path-shaped go test commands (fixture-relative).
//
// Phase 2 plan/contract only: may return multiple Cmds for nested modules.
// Production go-test execution still uses ModeWorkspaceSuite / ModeHubSuite.
//
// Rules:
//   - path/... → go test ./path/... in the module that owns path
//   - never expand to parent siblings (not go test ./tree/... when user said ./tree/mid/...)
//   - each nested go.mod under the path base → extra (cd nestedRoot && go test ./...)
//   - path without ... → go test ./path (single package pattern, no recursion)
func TranslatePath(userArg string, layout Layout) ([]Cmd, error) {
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

	owner := owningModule(base, mods)
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

	if recursive {
		prefix := base
		for _, m := range mods {
			if m == "." || m == owner {
				continue
			}
			if prefix == "." {
				out = append(out, Cmd{Dir: m, Pattern: "./..."})
				continue
			}
			if m == prefix || strings.HasPrefix(m, prefix+"/") {
				out = append(out, Cmd{Dir: m, Pattern: "./..."})
			}
		}
	}

	return dedupeCmds(out), nil
}

// NeedsPathShaped reports whether userArg + layout should use ModePathShaped
// instead of ModeWorkspaceSuite / ModeHubSuite.
//
// Phase 2 (plan/contract only): production still runs ModeWorkspaceSuite /
// ModeHubSuite via Plan; this helper is not wired into the workspace runner.
//
// True when the path selection crosses nested go.mod roots under the base
// (policy B/C multi-cmd). Mid-branch suite filtering stays in discovery
// (SubDir); default go test remains workspace/hub unless nested modules
// under the path require extra (cd nested && go test ./...) cmds.
//
// Optional TreeRoots (if set later) can mark mid-under-tree without go.mod.
func NeedsPathShaped(userArg string, layout Layout) bool {
	arg := strings.TrimSpace(userArg)
	if arg == "" || arg == "..." {
		return false
	}
	recursive := strings.HasSuffix(arg, "/...")
	if !recursive {
		// Non-recursive mid leaf still uses today's suite (one leaf filtered).
		// Path-shaped only when the path itself sits inside a nested module.
		base := strings.TrimPrefix(strings.TrimSuffix(arg, "/..."), "./")
		base = path.Clean(base)
		owner := owningModule(base, normalizeModules(layout.ModuleRoots))
		return owner != "."
	}
	base := strings.TrimSuffix(arg, "/...")
	base = strings.TrimPrefix(base, "./")
	if base == "" {
		base = "."
	}
	base = path.Clean(base)
	if base == "." {
		return hasNestedModule(layout)
	}
	return nestedModuleUnder(base, layout)
}

func hasNestedModule(layout Layout) bool {
	for _, m := range normalizeModules(layout.ModuleRoots) {
		if m != "." {
			return true
		}
	}
	return false
}

func nestedModuleUnder(prefix string, layout Layout) bool {
	prefix = path.Clean(prefix)
	for _, m := range normalizeModules(layout.ModuleRoots) {
		if m == "." {
			continue
		}
		if m == prefix || strings.HasPrefix(m, prefix+"/") {
			return true
		}
	}
	return false
}

// SuitePatternFromGen returns the package pattern for a workspace suite under genRoot.
func SuitePatternFromGen(genRoot, suiteAbsDir string) (string, error) {
	rel, err := filepath.Rel(genRoot, suiteAbsDir)
	if err != nil {
		return "", err
	}
	if rel == "." {
		return ".", nil
	}
	return "./" + filepath.ToSlash(rel), nil
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
		if strings.Count(out[i], "/") != strings.Count(out[j], "/") {
			return strings.Count(out[i], "/") > strings.Count(out[j], "/")
		}
		return out[i] < out[j]
	})
	return out
}

func owningModule(base string, mods []string) string {
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
