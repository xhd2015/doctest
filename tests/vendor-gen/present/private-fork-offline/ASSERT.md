## Expected

- `Run` succeeds.
- Gen `go.mod` is **offline-safe** for the private-fork modules.txt shape:
  1. **No bare require of the private fork** (`example.com/fork`) that forces
     network: must not have `require example.com/fork …` without a **filesystem**
     replace for that same path (`replace example.com/fork => <local path>`).
  2. **Prefer import path as require left-hand**:
     `require example.com/lib v1.0.0` (or equivalent require-block form).
  3. **Prefer filesystem replace to vendor content under the import path**:
     `replace example.com/lib => <modRoot>/vendor/example.com/lib`
     **or** an equivalent offline-safe graph that never leaves a bare private
     module path for tidy to download.
- Exactly one replace for the import-path left-hand when the offline form uses
  `example.com/lib` (no dual-replace with parent `lib => fork` **and** vendor).
- Other vendor modules (`example.com/nogo`) still get vendor require+replace.
- Project replace `replace example.com/app => <modRoot>` still present.

## Errors

- Classic TDD: **RED** today. `vendorBridgeForModRoot` rewrites require/replace
  left-hand to the non-local replacement path (`example.com/fork`) from
  modules.txt `A => B`, emitting e.g.:
  ```
  require example.com/fork v1.0.0   # BAD — private; tidy must fetch
  replace example.com/lib => example.com/fork v1.1.0  # parent, no FS target
  ```
  Desired offline-safe graph prefers:
  ```
  require example.com/lib v1.0.0
  replace example.com/lib => <modRoot>/vendor/example.com/lib
  ```
  (or equivalent that does not bare-require the private fork without a path replace).

```go
import (
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("WriteGoMod failed: %v", err)
	}
	goMod := resp.GoModContent
	if !strings.Contains(goMod, "module testcase") {
		t.Fatalf("expected module testcase, got:\n%s", goMod)
	}
	if !strings.Contains(goMod, "replace "+req.ModPath+" =>") {
		t.Fatalf("expected project replace, got:\n%s", goMod)
	}

	importPath := req.SampleModPath // example.com/lib
	importVer := req.SampleModVersion
	forkPath := req.ParentLocalRel   // example.com/fork
	// forkVer in ParentLocalAbs is fixture metadata; not required in gen graph.
	if importPath == "" || forkPath == "" {
		t.Fatal("Setup must set SampleModPath (import A) and ParentLocalRel (fork B)")
	}

	// --- 1. No bare require of private fork without filesystem replace ---
	if hasAnyRequire(goMod, forkPath) && !hasFilesystemReplace(goMod, forkPath) {
		t.Fatalf("offline-unsafe: bare require of private fork %q without filesystem replace "+
			"(forces network / go mod tidy download). gen go.mod:\n%s", forkPath, goMod)
	}

	// --- 2. Prefer require import path (A), not only private fork (B) ---
	if !hasRequire(goMod, importPath, importVer) {
		// Soft preference still fails the leaf: offline-safe form needs import path.
		t.Fatalf("expected require %s %s (import path as require left-hand for private-fork shape), got:\n%s",
			importPath, importVer, goMod)
	}

	// --- 3. Prefer filesystem replace of import path -> vendor/import path ---
	// Accept equivalent offline-safe: replace importPath => …/vendor/importPath
	// OR (if someone requires fork) filesystem replace of fork that points at
	// local content — but packages live under vendor/importPath, so preferred
	// is replace importPath => vendor/importPath.
	offlineOK := hasReplaceToVendor(goMod, importPath, req.ModRoot)
	if !offlineOK && hasFilesystemReplace(goMod, forkPath) {
		// Equivalent only if fork FS replace points at vendor content for import path.
		rhs := firstFilesystemReplaceRHS(goMod, forkPath)
		slash := filepath.ToSlash(rhs)
		if strings.Contains(slash, "/vendor/"+importPath) ||
			strings.HasSuffix(slash, "/vendor/"+importPath) {
			offlineOK = true
		}
	}
	if !offlineOK {
		t.Fatalf("expected offline-safe filesystem replace for private-fork shape: "+
			"prefer replace %s => %s/vendor/%s (or equivalent); got:\n%s",
			importPath, req.ModRoot, importPath, goMod)
	}

	// Dual-replace still forbidden for the import-path left-hand when that is
	// the offline form: at most one replace line for importPath.
	if n := countReplaceLHS(goMod, importPath); n > 1 {
		t.Fatalf("dual-replace forbidden: expected at most one replace for %s, got %d:\n%s",
			importPath, n, goMod)
	}

	// --- Other vendor module still injected ---
	if req.NoGoModPath != "" {
		if !hasRequire(goMod, req.NoGoModPath, req.NoGoModVersion) {
			t.Fatalf("expected vendor require %s %s for non-fork module, got:\n%s",
				req.NoGoModPath, req.NoGoModVersion, goMod)
		}
		if !hasReplaceToVendor(goMod, req.NoGoModPath, req.ModRoot) {
			t.Fatalf("expected vendor replace for %s, got:\n%s", req.NoGoModPath, goMod)
		}
	}
}

// hasAnyRequire reports any require line (single or block) for modPath, any version.
func hasAnyRequire(goMod, modPath string) bool {
	for _, line := range strings.Split(goMod, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "require "+modPath || strings.HasPrefix(trimmed, "require "+modPath+" ") {
			return true
		}
		// block form: "modPath version"
		fields := strings.Fields(trimmed)
		if len(fields) >= 1 && fields[0] == modPath {
			// avoid matching replace lines
			if !strings.HasPrefix(trimmed, "replace ") && !strings.HasPrefix(trimmed, "module ") {
				return true
			}
		}
	}
	return false
}

// hasFilesystemReplace reports replace modPath => <filesystem path> (not module+ver).
func hasFilesystemReplace(goMod, modPath string) bool {
	for _, line := range strings.Split(goMod, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "replace ") {
			continue
		}
		rest := strings.TrimSpace(strings.TrimPrefix(line, "replace "))
		parts := strings.Fields(rest)
		arrow := -1
		for i, p := range parts {
			if p == "=>" {
				arrow = i
				break
			}
		}
		if arrow < 1 || arrow+1 >= len(parts) {
			continue
		}
		if parts[0] != modPath {
			continue
		}
		right := parts[arrow+1:]
		// Filesystem: single token path-like (./ ../ abs, or containing /vendor/).
		if len(right) == 1 {
			r := right[0]
			if strings.HasPrefix(r, "./") || strings.HasPrefix(r, "../") ||
				filepath.IsAbs(r) || strings.Contains(r, string(filepath.Separator)) ||
				strings.Contains(filepath.ToSlash(r), "/") {
				// Module paths also contain / — treat as FS if looks local:
				// abs, relative, or contains "vendor" segment used by WriteGoMod.
				if filepath.IsAbs(r) || strings.HasPrefix(r, ".") ||
					strings.Contains(filepath.ToSlash(r), "/vendor/") {
					return true
				}
			}
		}
	}
	return false
}

func firstFilesystemReplaceRHS(goMod, modPath string) string {
	for _, line := range strings.Split(goMod, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "replace ") {
			continue
		}
		rest := strings.TrimSpace(strings.TrimPrefix(line, "replace "))
		parts := strings.Fields(rest)
		arrow := -1
		for i, p := range parts {
			if p == "=>" {
				arrow = i
				break
			}
		}
		if arrow < 1 || arrow+1 >= len(parts) || parts[0] != modPath {
			continue
		}
		right := parts[arrow+1:]
		if len(right) == 1 {
			r := right[0]
			if filepath.IsAbs(r) || strings.HasPrefix(r, ".") ||
				strings.Contains(filepath.ToSlash(r), "/vendor/") {
				return r
			}
		}
	}
	return ""
}

func countReplaceLHS(goMod, modPath string) int {
	n := 0
	for _, line := range strings.Split(goMod, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "replace ") {
			continue
		}
		rest := strings.TrimSpace(strings.TrimPrefix(line, "replace "))
		parts := strings.Fields(rest)
		if len(parts) >= 1 && parts[0] == modPath {
			n++
		}
	}
	return n
}
```
