# Scenario

**Feature**: gen-root source-input early-out for WriteGoMod / WriteGoModWithVendorBridges

```
# cold
WriteGoModWithVendorBridges(genDir, modRoot, …)
  -> gen go.mod (+ optional go.sum, overlay placeholders)
  -> doctest.gomod-src  (policy gomod-src=1 + source hashes)
  -> doctest.vendor-bridges.json

# warm hit (identical sources + integrity OK)
second WriteGoModWithVendorBridges
  -> no rewrite (mtime stable); tidy-done retained; bridges from cache

# miss (source or integrity)
second call rebuilds; tidy-done dropped iff gen mod/sum wrote
```

## Preconditions

- Package `github.com/xhd2015/doctest/libdoc/core` is importable.
- Each leaf uses isolated temp `ModRoot` and `GenDir` (`t.TempDir`); never
  shared mapping-gen.
- Coverage backfill: product already implements gomod-src cache — GREEN expected.
- Parallel-safe: no `t.Setenv`, `os.Chdir`, or process globals; snapshots on `req`.
- Do **not** read `DOCTEST_SESSION_ID` via `os.Getenv`.

## Steps

1. Root Setup sets default `ModPath` / `HasMod`.
2. Group/leaf Setup prepares fixtures, optional first write, and snapshots.
3. `Run` performs the measured `WriteGoModWithVendorBridges`; Assert inspects
   files, mtimes, tidy-done, and returned bridges.

## Context

- Artifact names: `doctest.gomod-src`, `doctest.vendor-bridges.json`,
  `doctest.tidy-done`, `doctest.gomod-fp` (must stay absent),
  `vendor-gomod-overlay.json`.
- Helpers: `prepareFreshGen`, `seedParentMod`, `seedVendorNogo`, `firstWrite`,
  `seedTidyDone`, `forceOldMtime`, `fillResponse`, `fileExists`, `readFileOrEmpty`.

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
	gomodSrcName      = "doctest.gomod-src"
	bridgesCacheName  = "doctest.vendor-bridges.json"
	tidyDoneName      = "doctest.tidy-done"
	gomodFpName       = "doctest.gomod-fp"
	overlayJSONName   = "vendor-gomod-overlay.json"
	defaultModPath    = "example.com/app"
	defaultNogoPath   = "example.com/nogo"
	defaultNogoVer    = "v0.1.0"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.ModPath == "" {
		req.ModPath = defaultModPath
	}
	if !req.HasMod {
		req.HasMod = true
	}
	if req.NogoModPath == "" {
		req.NogoModPath = defaultNogoPath
	}
	return nil
}

func prepareFreshGen(t *testing.T, req *Request, parentGoMod string) {
	t.Helper()
	req.ModRoot = t.TempDir()
	req.GenDir = filepath.Join(t.TempDir(), "gen")
	seedParentMod(t, req.ModRoot, req.ModPath, parentGoMod)
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

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// seedVendorNogo writes vendor/modules.txt + a module package without go.mod.
func seedVendorNogo(t *testing.T, req *Request) {
	t.Helper()
	if req.ModRoot == "" {
		t.Fatal("seedVendorNogo: ModRoot empty")
	}
	if req.NogoModPath == "" {
		req.NogoModPath = defaultNogoPath
	}
	vendor := filepath.Join(req.ModRoot, "vendor")
	modulesTxt := "# " + req.NogoModPath + " " + defaultNogoVer + "\n## explicit; go 1.17\n" + req.NogoModPath + "\n"
	writeFile(t, filepath.Join(vendor, "modules.txt"), modulesTxt)
	nogoDir := filepath.Join(vendor, filepath.FromSlash(req.NogoModPath))
	writeFile(t, filepath.Join(nogoDir, "nogo.go"), "package nogo\n\nconst X = 1\n")
	req.SeedVendorNogo = true
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func readFileOrEmpty(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

func forceOldMtime(t *testing.T, path string) time.Time {
	t.Helper()
	old := time.Now().Add(-3 * time.Hour).UTC().Truncate(time.Second)
	if err := os.Chtimes(path, old, old); err != nil {
		t.Fatalf("Chtimes %s: %v", path, err)
	}
	st, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat after Chtimes %s: %v", path, err)
	}
	return st.ModTime()
}

func seedTidyDone(t *testing.T, genDir string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(genDir, tidyDoneName), []byte("done"), 0644); err != nil {
		t.Fatalf("seed tidy-done: %v", err)
	}
}

func firstWrite(t *testing.T, req *Request) []core.VendorBridgeMapping {
	t.Helper()
	bridges, err := core.WriteGoModWithVendorBridges(req.GenDir, req.ModRoot, req.ModPath, req.HasMod,
		req.WithAssertReplace, req.AssertCacheDir,
		req.WithSessionReplace, req.SessionCacheDir)
	if err != nil {
		t.Fatalf("first WriteGoModWithVendorBridges: %v", err)
	}
	return bridges
}

func snapshotGoModMtime(t *testing.T, req *Request) {
	t.Helper()
	p := filepath.Join(req.GenDir, "go.mod")
	req.SnapGoModMtimeBefore = forceOldMtime(t, p)
	req.SnapGoModContentBefore = readFileOrEmpty(p)
}

func snapshotGomodSrc(t *testing.T, req *Request) {
	t.Helper()
	p := filepath.Join(req.GenDir, gomodSrcName)
	req.SnapGomodSrcBefore = readFileOrEmpty(p)
	if req.SnapGomodSrcBefore == "" {
		t.Fatalf("expected %s after first write", gomodSrcName)
	}
	req.SnapGomodSrcMtime = forceOldMtime(t, p)
}

func snapshotBridges(t *testing.T, req *Request, bridges []core.VendorBridgeMapping) {
	t.Helper()
	req.SnapBridgeCount = len(bridges)
	if len(bridges) == 0 {
		return
	}
	req.SnapBridgeRoot = bridges[0].BridgeRoot
	req.SnapBridgeModulePath = bridges[0].ModulePath
	req.SnapPlaceholderPath = bridges[0].BridgeRoot
	if req.SnapPlaceholderPath != "" && fileExists(req.SnapPlaceholderPath) {
		req.SnapPlaceholderMtime = forceOldMtime(t, req.SnapPlaceholderPath)
	}
}

func fillResponse(t *testing.T, req *Request, resp *Response, bridges []core.VendorBridgeMapping) {
	t.Helper()
	modFile := filepath.Join(req.GenDir, "go.mod")
	srcFile := filepath.Join(req.GenDir, gomodSrcName)
	brFile := filepath.Join(req.GenDir, bridgesCacheName)
	overlayFile := filepath.Join(req.GenDir, overlayJSONName)

	resp.GoModContent = readFileOrEmpty(modFile)
	resp.GomodSrcContent = readFileOrEmpty(srcFile)
	resp.GomodSrcExists = fileExists(srcFile)
	resp.BridgesJSONExists = fileExists(brFile)
	resp.BridgesJSONContent = readFileOrEmpty(brFile)
	resp.GomodFpExists = fileExists(filepath.Join(req.GenDir, gomodFpName))
	resp.TidyDoneExists = fileExists(filepath.Join(req.GenDir, tidyDoneName))
	resp.OverlayJSONExists = fileExists(overlayFile)

	resp.GoModMtimeBefore = req.SnapGoModMtimeBefore
	resp.GomodSrcMtimeBefore = req.SnapGomodSrcMtime
	resp.PlaceholderMtimeBefore = req.SnapPlaceholderMtime
	resp.PlaceholderPath = req.SnapPlaceholderPath

	if st, err := os.Stat(modFile); err == nil {
		resp.GoModMtimeAfter = st.ModTime()
	}
	if st, err := os.Stat(srcFile); err == nil {
		resp.GomodSrcMtimeAfter = st.ModTime()
	}

	resp.BridgeCount = len(bridges)
	resp.BridgeRoots = make([]string, 0, len(bridges))
	resp.BridgeModulePaths = make([]string, 0, len(bridges))
	for _, b := range bridges {
		resp.BridgeRoots = append(resp.BridgeRoots, b.BridgeRoot)
		resp.BridgeModulePaths = append(resp.BridgeModulePaths, b.ModulePath)
		if resp.PlaceholderPath == "" && b.BridgeRoot != "" {
			resp.PlaceholderPath = b.BridgeRoot
		}
	}
	if resp.PlaceholderPath != "" {
		resp.PlaceholderExists = fileExists(resp.PlaceholderPath)
		if st, err := os.Stat(resp.PlaceholderPath); err == nil {
			resp.PlaceholderMtimeAfter = st.ModTime()
		}
	}
}

func requireSrcAndBridges(t *testing.T, resp *Response) {
	t.Helper()
	if !resp.GomodSrcExists {
		t.Fatalf("expected %s", gomodSrcName)
	}
	if !resp.BridgesJSONExists {
		t.Fatalf("expected %s", bridgesCacheName)
	}
}

func fingerprintHasPolicyVersion(content string) bool {
	return strings.HasPrefix(strings.TrimSpace(content), "version gomod-src=1")
}
```
