# Scenario

**Feature**: generate writers share one content-hash manifest under each gen root

```
# first write into gen root
WriteGoMod / WriteIfChanged
  -> write targets when needed
  -> doctest.gen-manifest[path] = content-hash
  -> flush manifest if map changed (content-stable)

# warm identical
desired hash == manifest[path]
  -> skip target ReadFile + skip rewrite (mtime stable)
  -> tidy-done retained when go.mod/go.sum did not write

# no legacy skip file
doctest.gomod-fp never written
```

## Preconditions

- Package `github.com/xhd2015/doctest/libdoc/core` is importable.
- Each leaf uses an isolated temp gen root (`t.TempDir`); never the shared
  mapping-gen cache.
- Classic TDD: expect **RED** until unified `doctest.gen-manifest` replaces
  `doctest.gomod-fp` and WriteIfChanged consults the manifest.
- Do **not** read `DOCTEST_SESSION_ID` via `os.Getenv`; session identity is
  unused for these library leaves.

## Steps

1. Descendant Setup configures `req` (Mode, ModPath, flags) and prepares
   `ModRoot` / `GenDir` fixtures.
2. Root helpers seed parent `go.mod`, force old mtimes, and snapshot artifacts
   for warm/change scenarios.
3. `Run` performs the measured write; Assert inspects gen-root side effects.

## Context

- Manifest file name under test: `doctest.gen-manifest` (constant name may live
  in core; tests look for this path at gen root).
- Helpers: `seedParentMod`, `forceOldMtime`, `manifestPath`, `fillResponse`,
  `writeGenRelFile` (WriteIfChanged-facing write of a relative gen path —
  preferably via core’s unified writer once exported; falls back to
  `WriteFormattedGo` under gen root for Go sources).
- `runKind` is not used; `req.Mode` selects the Run branch.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/core"
)

const (
	genManifestName = "doctest.gen-manifest"
	gomodFpName     = "doctest.gomod-fp"
	tidyDoneName    = "doctest.tidy-done"
)

// Package snapshots: Setup records "before" state for warm/second-call leaves;
// fillResponse copies them into Response for Assert.
var (
	snapGoModMtimeBefore    time.Time
	snapManifestMtimeBefore time.Time
	snapTargetMtimeBefore   time.Time
	snapManifestEntryBefore string
	snapGoModContentBefore  string
	snapManifestContentBefore string
)

func resetSnapshots() {
	snapGoModMtimeBefore = time.Time{}
	snapManifestMtimeBefore = time.Time{}
	snapTargetMtimeBefore = time.Time{}
	snapManifestEntryBefore = ""
	snapGoModContentBefore = ""
	snapManifestContentBefore = ""
}

func Setup(t *testing.T, req *Request) error {
	resetSnapshots()
	// Defaults; leaves override Mode / paths / flags.
	if req.ModPath == "" {
		req.ModPath = "example.com/app"
	}
	if !req.HasMod {
		req.HasMod = true
	}
	return nil
}

func manifestPath(genDir string) string {
	return filepath.Join(genDir, genManifestName)
}

func seedParentMod(t *testing.T, modRoot, modPath, body string) {
	t.Helper()
	if body == "" {
		body = "module " + modPath + "\n\ngo 1.21\n"
	}
	if err := os.MkdirAll(modRoot, 0755); err != nil {
		t.Fatalf("mkdir modRoot: %v", err)
	}
	if err := os.WriteFile(filepath.Join(modRoot, "go.mod"), []byte(body), 0644); err != nil {
		t.Fatalf("write parent go.mod: %v", err)
	}
}

func prepareFreshGen(t *testing.T, req *Request, parentGoMod string) {
	t.Helper()
	req.ModRoot = t.TempDir()
	req.GenDir = filepath.Join(t.TempDir(), "gen")
	seedParentMod(t, req.ModRoot, req.ModPath, parentGoMod)
}

func forceOldMtime(t *testing.T, path string) time.Time {
	t.Helper()
	old := time.Now().Add(-2 * time.Hour).UTC().Truncate(time.Second)
	if err := os.Chtimes(path, old, old); err != nil {
		t.Fatalf("Chtimes %s: %v", path, err)
	}
	st, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat after Chtimes %s: %v", path, err)
	}
	return st.ModTime()
}

func readFileOrEmpty(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func seedTidyDone(t *testing.T, genDir string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(genDir, tidyDoneName), []byte("done"), 0644); err != nil {
		t.Fatalf("seed tidy-done: %v", err)
	}
}

// writeGenRelFile writes content at genRoot/rel using the unified gen write
// path when available. Prefer content-hash WriteIfChanged under gen root so
// the manifest is updated. Uses WriteFormattedGo for .go sources so post-format
// bytes match production writers.
func writeGenRelFile(t *testing.T, genRoot, rel, content string) error {
	t.Helper()
	if genRoot == "" || rel == "" {
		t.Fatalf("writeGenRelFile: genRoot and rel required")
	}
	rel = filepath.ToSlash(rel)
	abs := filepath.Join(genRoot, filepath.FromSlash(rel))
	if strings.HasSuffix(rel, ".go") {
		// WriteFormattedGo formats then content-stable writes; implementer
		// must also update doctest.gen-manifest[rel] (and flush if needed).
		if err := core.WriteFormattedGo(abs, content); err != nil {
			return err
		}
		// After implementation, writers that know genRoot should record the
		// relative path hash. Call optional TestExported hook if present so
		// library unit coverage can force manifest update without full generate.
		if fn := genManifestRecordHook(); fn != nil {
			return fn(genRoot, rel, abs)
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
		return err
	}
	// Non-Go: still expect unified WriteIfChanged semantics once exported.
	if fn := genWriteIfChangedHook(); fn != nil {
		return fn(genRoot, rel, []byte(content))
	}
	return os.WriteFile(abs, []byte(content), 0644)
}

// Optional hooks: implementer may set these via TestExported_ helpers so leaves
// stay RED until real WriteIfChanged + manifest record exist. Default nil.
var (
	testGenManifestRecord func(genRoot, rel, absPath string) error
	testGenWriteIfChanged func(genRoot, rel string, data []byte) error
)

func genManifestRecordHook() func(genRoot, rel, absPath string) error {
	return testGenManifestRecord
}

func genWriteIfChangedHook() func(genRoot, rel string, data []byte) error {
	return testGenWriteIfChanged
}

func fillResponse(t *testing.T, req *Request, resp *Response) {
	t.Helper()
	modFile := filepath.Join(req.GenDir, "go.mod")
	manFile := manifestPath(req.GenDir)
	resp.GoModContent = readFileOrEmpty(modFile)
	resp.ManifestContent = readFileOrEmpty(manFile)
	resp.ManifestExists = fileExists(manFile)
	resp.GomodFpExists = fileExists(filepath.Join(req.GenDir, gomodFpName))
	resp.TidyDoneExists = fileExists(filepath.Join(req.GenDir, tidyDoneName))
	resp.GoModMtimeBefore = snapGoModMtimeBefore
	resp.ManifestMtimeBefore = snapManifestMtimeBefore
	resp.TargetMtimeBefore = snapTargetMtimeBefore
	resp.ManifestEntryBefore = snapManifestEntryBefore
	if st, err := os.Stat(modFile); err == nil {
		resp.GoModMtimeAfter = st.ModTime()
	}
	if st, err := os.Stat(manFile); err == nil {
		resp.ManifestMtimeAfter = st.ModTime()
	}
	if req.RelPath != "" {
		abs := filepath.Join(req.GenDir, filepath.FromSlash(req.RelPath))
		resp.TargetContent = readFileOrEmpty(abs)
		if st, err := os.Stat(abs); err == nil {
			resp.TargetMtimeAfter = st.ModTime()
		}
		resp.ManifestEntryAfter = findManifestLine(resp.ManifestContent, req.RelPath)
	} else {
		resp.ManifestEntryAfter = findManifestLine(resp.ManifestContent, "go.mod")
	}
}

func findManifestLine(manifest, rel string) string {
	rel = filepath.ToSlash(rel)
	for _, line := range strings.Split(manifest, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Accept "version N" and "path hash" lines; match path token.
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		if fields[0] == rel || strings.HasPrefix(line, rel+" ") || strings.Contains(line, `"`+rel+`"`) {
			return line
		}
	}
	// Fallback: any line containing the relative path.
	for _, line := range strings.Split(manifest, "\n") {
		if strings.Contains(line, rel) {
			return strings.TrimSpace(line)
		}
	}
	return ""
}

func firstWriteGoMod(t *testing.T, req *Request) {
	t.Helper()
	if err := core.WriteGoMod(req.GenDir, req.ModRoot, req.ModPath, req.HasMod,
		req.WithAssertReplace, req.AssertCacheDir,
		req.WithSessionReplace, req.SessionCacheDir); err != nil {
		t.Fatalf("first WriteGoMod: %v", err)
	}
}

func snapshotGoModMtime(t *testing.T, req *Request) {
	t.Helper()
	p := filepath.Join(req.GenDir, "go.mod")
	snapGoModMtimeBefore = forceOldMtime(t, p)
	snapGoModContentBefore = readFileOrEmpty(p)
}

func snapshotManifestMtime(t *testing.T, req *Request) {
	t.Helper()
	p := manifestPath(req.GenDir)
	if !fileExists(p) {
		// Pre-implementation: leave zero; Assert fails on existence/mtime.
		return
	}
	snapManifestMtimeBefore = forceOldMtime(t, p)
	snapManifestContentBefore = readFileOrEmpty(p)
	snapManifestEntryBefore = findManifestLine(snapManifestContentBefore, "go.mod")
}

func snapshotTargetMtime(t *testing.T, req *Request) {
	t.Helper()
	if req.RelPath == "" {
		t.Fatalf("snapshotTargetMtime: RelPath required")
	}
	abs := filepath.Join(req.GenDir, filepath.FromSlash(req.RelPath))
	snapTargetMtimeBefore = forceOldMtime(t, abs)
	man := readFileOrEmpty(manifestPath(req.GenDir))
	snapManifestEntryBefore = findManifestLine(man, req.RelPath)
	snapManifestContentBefore = man
	if fileExists(manifestPath(req.GenDir)) {
		snapManifestMtimeBefore = forceOldMtime(t, manifestPath(req.GenDir))
	}
}

func requireUnifiedLayout(t *testing.T, resp *Response) {
	t.Helper()
	if !resp.ManifestExists {
		t.Fatalf("expected %s at gen root", genManifestName)
	}
	if resp.GomodFpExists {
		t.Fatalf("legacy %s must not exist under unified manifest", gomodFpName)
	}
}
```
